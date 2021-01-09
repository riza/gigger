package gigger

import (
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog/log"

	"github.com/riza/gigger/pkg/config"
	"github.com/riza/gigger/pkg/git"
	"github.com/riza/gigger/pkg/task"
)

type Pool struct {
	pool *ants.PoolWithFunc
	task *task.Task
	conf *config.Config
	Wg   sync.WaitGroup
	git  *git.Git
}

type URL struct {
	URL      string
	isIndex  bool
	isObject bool
	fileName string
}

func NewPool(conf *config.Config, t *task.Task) (*Pool, error) {
	p := &Pool{}
	p.Wg = sync.WaitGroup{}

	var err error
	p.pool, err = ants.NewPoolWithFunc(conf.Thread, p.process)
	if err != nil {
		return nil, err
	}

	p.git = git.NewGit()
	p.conf = conf
	p.task = t

	if string(p.conf.URL[len(p.conf.URL)-1]) != "/" {
		p.conf.URL += "/"
	}

	return p, nil
}

func (p *Pool) process(data interface{}) {
	url, check := data.(URL)
	if !check {
		p.Wg.Done()
	}

	status, body, _ := p.task.Client.Get(nil, url.URL)
	if status != 200 {
		log.Error().Msgf("[%d] %s", status, url.URL)
		p.Wg.Done()
		return
	} else {
		if p.conf.Verbose {
			log.Info().Msgf("[%d] %s", status, url.URL)
		}
	}

	if url.isIndex {
		index, err := p.git.ParseIndex(body)
		if err != nil {
			log.Error().Msgf("%s parsing error [%s]", url.URL, err.Error())
			p.Wg.Done()
			return
		}
		if p.conf.Verbose {
			log.Info().Msgf("%s parsed", url.URL)
		}

		for _, entry := range index.Entries {
			objectURL := URL{
				p.conf.URL + ".git/objects/" + string(entry.SHA1[0]) + string(entry.SHA1[1]) + "/" + entry.SHA1[2:],
				false,
				true,
				entry.Name,
			}

			if p.conf.Verbose {
				log.Info().Msgf("%s object found, downloading [%s]", objectURL.URL, objectURL.fileName)
			}

			p.pool.Invoke(objectURL)
			p.Wg.Add(1)
		}
	}
	if url.isObject {
		object, err := p.git.ParseObject(body)
		if err != nil {
			log.Error().Msgf("%s parsing error [%s]", url.URL, err.Error())
			p.Wg.Done()
			return
		}

		err = p.task.SaveFile(url.fileName, object)
		if err != nil {
			log.Error().Msgf("%s save file error [%s]", url.fileName, err.Error())
			p.Wg.Done()
			return
		}
		log.Info().Msgf("%s downloaded", url.fileName)
	}

	p.Wg.Done()

}

func (p *Pool) Run() error {
	for _, u := range p.generateList() {
		url := URL{u, u == p.conf.URL+".git/index", false, ""}
		err := p.pool.Invoke(url)
		if err != nil {
			return err
		}
		p.Wg.Add(1)
	}
	return nil
}

func (p *Pool) generateList() (list []string) {
	for path, typ := range GitFolderStructure {
		fileMap, isMap := typ.(map[string]bool)
		if isMap {
			for path2, _ := range fileMap {
				list = append(list, p.conf.URL+path+path2)
			}
			continue
		}
		list = append(list, p.conf.URL+path)
	}
	return
}

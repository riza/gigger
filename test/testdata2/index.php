<?php
#The source of this file is top secret, I hope no one will see it.
$file = $_GET['file'];

if(isset($file))
    include("pages/$file");
?>
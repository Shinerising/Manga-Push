<?php

require_once('func.php');
require_once('fpdf.php');
require_once('PHPMailer/PHPMailerAutoload.php');

header('Content-Type: application/json; charset=utf-8');

$des=$_GET['email'];
$from=$_GET['para1'];
$name=$_GET['para2'];
$mgid=$_GET['para3'];
$title=$_GET['title'];
$author=$_GET['author'];

$file = 'pdf/'.$name.$mgid.'.pdf';
$imgdir = 'manga/'.$name.$mgid.'/';
if(file_exists($file)){
	if(sendMail($des, $title, $file)) {
		pushMessage('succeed',$title,'Push Manga Successful!');
	}
	else {
		pushMessage('error',$title,'Push Manga Failed!');
	}
}
else if($from=="http://www.ishuhui.com/") {
	$images = loadSHUHUIImages($mgid);
	if(count($images)>0) {
		if(createPDF($imgdir, $file, $images, $title, $author)) {
			if(sendMail($des, $title, $file)) {
				pushMessage('succeed',$title,'Push Manga Successful!');
			}
			else {
				pushMessage('error',$title,'Push Manga Failed!');
			}
		}
		else {
			pushMessage('error',$title,'Create Manga File Failed!');
		}
	}
	else {
		pushMessage('error',$title,'Fetching Images Failed!');
	}
}
else if($from=="http://manhua.dmzj.com/") {
	$images = loadDMZJImages($name, $mgid);
	if(count($images)>0) {
		if(createPDF($imgdir, $file, $images, $title, $author)) {
			if(sendMail($des, $title, $file)) {
				pushMessage('succeed',$title,'Push Manga Successful!');
			}
			else {
				pushMessage('error',$title,'Push Manga Failed!');
			}
		}
		else {
			pushMessage('error',$title,'Create Manga File Failed!');
		}
	}
	else {
		pushMessage('error',$title,'Fetching Images Failed!');
	}
}
else {
	pushMessage('error',$title,'Invalid Parameters!');
}


?>
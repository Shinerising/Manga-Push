<?php

function inCategory($str, $arr){
	foreach($arr as $item) {
		if($item==(string)$str) return true;
	}
	return false;
}

if(file_exists('config_default.php')) require_once('config_default.php');
else require_once('config.php');
require_once('func.php');
require_once('fpdf.php');
require_once('PHPMailer/PHPMailerAutoload.php');

error_reporting(0);

header('Content-type: text/xml; charset=utf-8');
//header('Content-type: text/plain; charset=utf-8');

$xml=simplexml_load_file("settings.xml");

$cuser=$_GET['user'];

$email_list=$GLOBALS['email_list'];

$des = $email_list[$cuser];

$ruser=null;
foreach($xml->xpath('//users/user') as $user) {
	if($user['name']==$cuser) $ruser=$user;
}
$xmltitle='No New Manga';

if($ruser==null || time()-$ruser->lastvisit<$GLOBALS['push_interval']*60){
}
else {
	$items=$ruser->xpath('list/item');
	foreach($items as $item) {
		if($item->site=='dmzj.com') {
			$content=file_get_contents($item->feed);
			if($content=='') {
				$xmltitle='Get Images Feed Error!';
				break;
			}
			else {
				$feed=simplexml_load_string($content, null, LIBXML_NOCDATA);
				$manga_items=$feed->xpath('//channel/item');
				$request_key=-1;
				$id=0;
				foreach ($manga_items as $key=>$value) {
					preg_match("/http:\/\/manhua.dmzj.com\/.*?\/(.*?).shtml?/", $value->link, $output);
					$id = $output[1];
					if($id!=$item->lastlink) $request_key=$key;
					else break;
				}
				if($request_key==-1)continue;

				preg_match("/title='(.*?)'/", $manga_items[$request_key]->description, $output);
				$episode=$output[1];
				$title=$item->title.' '.$episode;
				$author=$item->author;

				preg_match("/http:\/\/manhua.dmzj.com\/.*?\/(.*?).shtml?/", $manga_items[$request_key]->link, $output);
				$mgid = $output[1];
				preg_match("/http:\/\/manhua.dmzj.com\/(.*?)\//", $manga_items[$request_key]->link, $output);
				$name = $output[1];

				$file = 'pdf/'.$name.$mgid.'.pdf';
				$imgdir = 'manga/'.$name.$mgid.'/';	
				
				if(file_exists($file)){
					if(sendMail($des, $title, $file)) {
						$xmltitle='Push Manga Successful!';
						$ruser->lastlink=$manga_items[$request_key]->link;
						$ruser->lastpubDate=$manga_items[$request_key]->pubDate;
						$ruser->lastdescription=$title;
						$item->lastlink=$mgid;
						$ruser->lastrequest=time();
						break;
					}
					else {
						$xmltitle='Push Manga Failed!';
						break;
					}
				}
				$images = loadDMZJImages($name, $mgid);
				if(count($images)>0) {
					if(createPDF($imgdir, $file, $images, (string)$title, (string)$author)) {
						if(sendMail($des, $title, $file)) {
							$xmltitle='Push Manga Successful!';
							$ruser->lastlink=$manga_items[$request_key]->link;
							$ruser->lastpubDate=$manga_items[$request_key]->pubDate;
							$ruser->lastdescription=$title;
							$ruser->lastrequest=time();
							$item->lastlink=$mgid;
							break;
						}
						else {
							$xmltitle='Push Manga Failed!';
							break;
						}
					}
					else {
						$xmltitle='Create Manga File Failed!';
						break;
					}
				}
				else {
					$xmltitle='Fetching Images Failed!';
					break;
				}
			}
		}
		else if($item->site=='ishuhui.com') {
			$content=file_get_contents($item->feed);
			if($content=='') {
				$xmltitle='Get Images Feed Error!';
				break;
			}
			else {
				$feed=simplexml_load_string($content, null, LIBXML_NOCDATA);
				$manga_items=$feed->xpath('//channel/item');
				$request_key=-1;
				$id=0;
				foreach ($manga_items as $key=>$value) {
					if(inCategory($item->title, $value->xpath('category'))) {
						preg_match("/http:\/\/www.ishuhui.com\/archives\/(.*)/", $value->link, $output);
						$id = $output[1];
						if($id!=$item->lastlink) $request_key=$key;
						else break;
					}
				}
				if($request_key==-1)continue;

				preg_match("/第(.*?)话/", $manga_items[$request_key]->title, $output);
				$episode=$output[0];
				$title=$item->title.' '.$episode;
				$author=$item->author;

				preg_match("/http:\/\/www.ishuhui.com\/archives\/(.*)/", $manga_items[$request_key]->link, $output);
				$mgid = $output[1];
				$name = $item->name;

				$file = 'pdf/'.$name.$mgid.'.pdf';
				$imgdir = 'manga/'.$name.$mgid.'/';

				if(file_exists($file)){
					if(sendMail($des, $title, $file)) {
						$xmltitle='Push Manga Successful!';
						$ruser->lastlink=$manga_items[$request_key]->link;
						$ruser->lastpubDate=$manga_items[$request_key]->pubDate;
						$ruser->lastdescription=$title;
						$item->lastlink=$mgid;
						$ruser->lastrequest=time();
						break;
					}
					else {
						$xmltitle='Push Manga Failed!';
						break;
					}
				}
				$images = loadSHUHUIImages($mgid);
				if(count($images)>0) {
					if(createPDF($imgdir, $file, $images, (string)$title, (string)$author)) {
						if(sendMail($des, $title, $file)) {
							$xmltitle='Push Manga Successful!';
							$ruser->lastlink=$manga_items[$request_key]->link;
							$ruser->lastpubDate=$manga_items[$request_key]->pubDate;
							$ruser->lastdescription=$title;
							$ruser->lastrequest=time();
							$item->lastlink=$mgid;
							break;
						}
						else {
							$xmltitle='Push Manga Failed!';
							break;
						}
					}
					else {
						$xmltitle='Create Manga File Failed!';
						break;
					}
				}
				else {
					$xmltitle='Fetching Images Failed!';
					break;
				}
			}

		}
	}

	$ruser->lastvisit=time();
    $xml->asXml('settings.xml');
}

?>
<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
	<channel>
		<title>Manga Push Feed</title>
		<description></description>
		<link><?php echo "http://$_SERVER[HTTP_HOST]$_SERVER[REQUEST_URI]"; ?></link>
		<ttl>60</ttl>
		<item>
    		<title><?php echo $xmltitle; ?></title>
    		<link><?php echo $ruser->lastlink; ?></link>
    		<pubDate><?php echo $ruser->lastpubDate; ?></pubDate>
    		<description><?php echo $ruser->lastdescription; ?></description>
    		<guid><?php echo $ruser->lastlink; ?></guid>
		</item>
	</channel>
</rss>
<?php

if(file_exists('config_default.php')) require_once('config_default.php');
else require_once('config.php');
require_once('fpdf.php');
require_once('PHPMailer/PHPMailerAutoload.php');

//  Print JSON strings to report pushing status
function pushMessage($status, $title, $msg) {
    $data['status']=$status;
    $data['title']=$title;
    $data['msg']=$msg;
    echo json_encode($data);
    exit();
}

//  Check Image Ratio
//  Return the Ratio of Width & Height
function checkImage($image, $format) {

    if ($format=='.jpg')
        $imageTmp=imagecreatefromjpeg($image);
    else if ($format=='.png')
        $imageTmp=imagecreatefrompng($image);
    else if ($format=='.gif')
        $imageTmp=imagecreatefromgif($image);
    else if ($format=='.bmp')
        $imageTmp=imagecreatefromwbmp($image);
    else return 0;

    $width = imagesx($imageTmp);
    $height = imagesy($imageTmp);
    imagedestroy($imageTmp);

    return $width / $height;
}

// Compress Image
// Return the Ratio of Width & Height
function compressImage($image, $desimage, $format, $maxH, $quality) {

    if ($format=='.jpg')
        $imageTmp=imagecreatefromjpeg($image);
    else if ($format=='.png')
        $imageTmp=imagecreatefrompng($image);
    else if ($format=='.gif')
        $imageTmp=imagecreatefromgif($image);
    else if ($format=='.bmp')
        $imageTmp=imagecreatefromwbmp($image);
    else return 0;

    $width = imagesx($imageTmp);
    $height = imagesy($imageTmp);

    if($maxH!=-1){
        $thumb_width = round($width * $maxH / $height);
        $thumb_height = round($maxH );
        $thumb = imagecreatetruecolor($thumb_width, $thumb_height);
        imagecopyresampled($thumb, $imageTmp, 0, 0, 0, 0, $thumb_width, $thumb_height, $width, $height);
        imagejpeg($thumb, $desimage, $quality);
        imagedestroy($thumb);
    } else {
        imagejpeg($imageTmp, $desimage, $quality);
    }
    
    imagedestroy($imageTmp);

    return $width / $height;
}

//  Create PDF File
// @param $dir: Images Directory
// @param $file: Output PDF File Path
// @param $images: Images URL array
// @param $name: PDF Title
// @param $author: Author
// @return If succeed return true
function createPDF($dir, $file, $images, $name, $author) {
    $pdf = new FPDF('P','mm',$GLOBALS['page_size']);
    $pdf->setTitle($name, true);
    $pdf->setSubject($name, true);
    $pdf->setAuthor($author, true);

    if(!is_dir('manga/')) mkdir('manga/');
    if(!is_dir($dir)) mkdir($dir);

    $imgcount=count($images);
    $quality=20;
    $maxH=480;
    if($imgcount<16){
        $quality = 100;
        $maxH = -1;
    } 
    else if($imgcount<32){
        $quality = 80;
        $maxH = -1;
    } 
    else if($imgcount<64){
        $quality = 60;
        $maxH = 1080;
    } 
    else if($imgcount<128){
        $quality = 40;
        $maxH = 960;
    } 
    if($imgcount>200) return false;

    foreach ($images as $key=>$value) {
        $exploded = explode('.',$value);
        $ext = $exploded[count($exploded) - 1]; 

        if (preg_match('/jpg|jpeg/i',$ext))
            $format='.jpg';
        else if (preg_match('/png/i',$ext))
            $format='.png';
        else if (preg_match('/gif/i',$ext))
            $format='.gif';
        else if (preg_match('/bmp/i',$ext))
            $format='.bmp';
        else $format='.jpg';

        if($format=='.gif') continue;

        $ratio=1;
        $imagefile=$dir.'image'.$key.'.jpg';
        if(file_exists($imagefile)){
            $ratio = checkImage($imagefile, '.jpg');
        }
        else {
            $imagefile=$dir.'image'.$key.$format;
            if(!copy($value, $imagefile)) return false;
            $ratio = compressImage($imagefile, $dir.'image'.$key.'.jpg', $format, $maxH, $quality);
            $imagefile=$dir.'image'.$key.'.jpg';
        }
        $pdf->AddPage();
        if($ratio > 1) {
            $pdf->Image($imagefile,92 - $ratio * 122,0,0,122);
            $pdf->AddPage();
            $pdf->Image($imagefile,0,0,0,122);

        }
        else if($ratio > 92 / 122) {
            $pdf->Image($imagefile,0,0,92,0);
        }
        else {
            $pdf->Image($imagefile,(92 - 122 * $ratio) / 2,0,0,122);
        }
    }

    $dir='pdf/';
    if(!is_dir($dir)) mkdir($dir);
    $pdf->Output('F', $file);
    return true;
}

//  Send Email to Kindle Cloud Server
//  @param $des: The Kindle Email Address
//  @param $title: Manga Title
//  @param $file: PDF File Path
//  @return If succeed return true
function sendMail($des, $title, $file) {
    $mail = new PHPMailer;
    $mail->CharSet = 'UTF-8';
    $mail->isSMTP();
    $mail->Host = $GLOBALS['mail_host'];
    $mail->SMTPAuth = true;
    $mail->Username = $GLOBALS['mail_username'];
    $mail->Password = $GLOBALS['mail_password'];
    $mail->SMTPSecure = $GLOBALS['mail_SMTPSecure'];
    $mail->Port = $GLOBALS['mail_port'];

    $mail->setFrom($GLOBALS['mail_address'], $GLOBALS['mail_name']);
    $mail->addAddress($des); 

    $mail->addAttachment($file, $title.'.pdf', 'base64', 'application/pdf');
    $mail->isHTML(true);

    $mail->Subject = '[MangaPush]'.$title;
    $mail->Body    = 'Email from Manga Push';

    if(!$mail->send()) {
        return false;
    } else {
        return true;
    }
}

// Get the Images URL of DMZJ Manga
function loadDMZJImages($name, $id) {
    $url='http://m.dmzj.com/view/'.$name.'/'.$id.'.html';
    $content=file_get_contents($url);
    preg_match_all("/\{(.*?)\}/", $content, $output);
    $json=$output[0][0];
    preg_match_all("/http(.*?)(png|jpg|jpeg)/", $json, $output);
    $images=$output[0];
    foreach ($images as $key=>$value) {
        $images[$key] = json_decode('"'.$value.'"');
    }
    return $images;
}

// Get the Images URL of ISHUHUI Manga
function loadSHUHUIImages($id) {
    $url='http://www.ishuhui.com/post/'.$id;
    $content=file_get_contents($url);
    preg_match_all("/<img src=\"(.*?)\" alt/", $content, $output);
    $images=$output[1];
    array_shift($images);
    return $images;
}

?>
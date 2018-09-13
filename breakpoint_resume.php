<?php
/**
 * Created by PhpStorm.
 * User: lrj@mxj
 * Date: 2018/9/13
 * Time: 6:09 PM
 */

$host = "127.0.0.1:8088";
$file = "/Users/mxj/Desktop/微服务：从设计到部署.pdf";
$url = "http://" . $host . "/api/upload";
$session_id = "lrj";

$upload_size_per_request = 50000;

$file_size = filesize($file);
if ($upload_size_per_request > $file_size) {
    $upload_size_per_request = $file_size;
}

$run_num = intval(($file_size - 1) / $upload_size_per_request) + 1;

c('file_size:' . $file_size);
c("run_num:" . $run_num);

$file_content = file_get_contents($file);// Todo 换种方式读文件.

$start = 0;
$end = $upload_size_per_request - 1;
for ($i = 1; $i <= $run_num; $i++) {
    c("第" . $i . "次上传");
    $length = $i < $run_num ? $upload_size_per_request : $file_size - ($i - 1) * $upload_size_per_request;

    $header =
        "Host: " . $host . "\r\n" .
        "Content-Length: " . $length . "\r\n" .
        "Content-Disposition: attachment; filename=\"" . basename($file) . "\"\r\n" .
        "X-Content-Range: bytes " . $start . "-" . $end . "/" . $file_size . "\r\n" .
        "Session-ID: " . $session_id . "\r\n" .
        "Content-Type: application/octet-stream\r\n\r\n";
    $content = substr($file_content, $start, $length);

    $context = stream_context_create(array(
        'http' => array(
            'method' => 'POST',
            'header' => $header,
            'content' => $content,
            'protocol_version' => "1.1"
        )
    ));

    $result = file_get_contents($url, false, $context);
    // Todo 检查$result正确性 $http_response_header.
    $start = $end + 1;
    $end = ($i == $run_num - 1) ? $file_size - 1 : ($start + $length - 1);
}


function c($str)
{
    echo $str;
    echo "\r\n";
}

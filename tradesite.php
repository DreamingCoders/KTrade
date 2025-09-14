<?php
//bar

/*
// create a new cURL resource
$ch = curl_init();
// set URL and other appropriate options
curl_setopt($ch, CURLOPT_URL, "http://www.example.net/api/foobar.php");
curl_setopt($ch, CURLOPT_HEADER, 0);
// grab URL and pass it to the browser
curl_exec($ch);
// close cURL resource, and free up system resources
curl_close($ch);
*/

$get_data = callAPI('GET', 'https://example.net/api/foobar.php, false);
$response = json_decode($get_data, true);
$errors = $response['response']['errors'];
$data = $response['response']['data'][0];
?>
<?php
//bar

/*
// create a new cURL resource
$ch = curl_init();
// set URL and other appropriate options
curl_setopt($ch, CURLOPT_URL, "http://www.example.net/api/foobar.php");
curl_setopt($ch, CURLOPT_HEADER, 0);
// grab URL and pass it to the browser
curl_exec($ch);
// close cURL resource, and free up system resources
curl_close($ch);
*/

$get_data = callAPI('GET', 'https://example.net/api/foobar.php, false);
$response = json_decode($get_data, true);
$errors = $response['response']['errors'];
$data = $response['response']['data'][0];
?>
<!DOCTYPE html>
<html>
<head>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.2.1/jquery.min.js"></script>
<script>
$(document).ready(function(){
   
   $('#content').load("https://www.example.net/api/foobar.php");
   $('#api').load("https://www.example.net/api/foobar.php", { "players[]": [ "online" ] });
});
$.getJSON( "https://www.example.net/api/test.json", function( data ) {
  var items = [];
  $.each( data, function( key, val ) {
    items.push( "<li id='" + key + "'>" + val + "</li>" );
  });
 
  $( "<ul/>", {
    "class": "my-new-list",
    html: items.join( "" )
  }).appendTo( "body" );
});
</script>
</head>
<body>

<h2>BPR API</h2>
<hr>

<p>JSON Data-</p>
<div id="content">Loading Content...</div>
<div id="api">Loading API...</div>


</body>
</html>
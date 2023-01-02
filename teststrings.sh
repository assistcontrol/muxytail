#!/bin/sh

cat <<EOF
{"level":"info","ts":1664590340.2631304,"logger":"http.log.access.log1","msg":"handled request","request":{"remote_ip":"104.225.5.94","remote_port":"53454","proto":"HTTP/2.0","method":"GET","host":"abg.ninja","uri":"/abg","headers":{"Sec-Ch-Ua":["\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\""],"Sec-Ch-Ua-Platform":["\"Windows\""],"Sec-Fetch-User":["?1"],"Sec-Fetch-Dest":["document"],"Sec-Ch-Ua-Mobile":["?0"],"Upgrade-Insecure-Requests":["1"],"User-Agent":["Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"],"Accept":["text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"],"Sec-Fetch-Site":["same-origin"],"Sec-Fetch-Mode":["navigate"],"Referer":["https://abg.ninja/abg"],"Accept-Encoding":["gzip, deflate, br"],"Accept-Language":["en-US,en;q=0.9"]},"tls":{"resumed":true,"version":772,"cipher_suite":4865,"proto":"h2","server_name":"abg.ninja"}},"user_id":"","duration":0.009557991,"size":3380,"status":200,"resp_headers":{"Referrer-Policy":["same-origin"],"Strict-Transport-Security":["max-age=31536000; includeSubdomains; preload"],"X-Content-Type-Options":["nosniff"],"X-Xss-Protection":["1; mode=block"],"Alt-Svc":["h3=\":443\"; ma=2592000"],"Content-Encoding":["gzip"],"Content-Length":["3380"],"X-Frame-Options":["sameorigin"],"Date":["Sat, 01 Oct 2022 02:12:20 GMT"],"Vary":["Accept-Encoding"],"Content-Type":["text/html;charset=UTF-8"],"Content-Security-Policy":["frame-ancestors 'self'"]}}
EOF

echo '1: Color 1 string'
echo '2: Color 2 string'
echo '3: Color 3 string'
echo '4: Danger string'
echo '5: Uncolored string'

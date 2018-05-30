var submitFeedback = function(){
    var form = document.getElementById('feedback_form');
    var body = new Object();
    body.title = form.title.value;
    body.email = form.email.value;
    body.content = form.content.value;

    var httpRequest;
    if (window.XMLHttpRequest) { // 모질라, 사파리등 그외 브라우저, ...
        httpRequest = new XMLHttpRequest();
    } else if (window.ActiveXObject) { // IE 8 이상
        httpRequest = new ActiveXObject("Microsoft.XMLHTTP");
    }
    httpRequest.onreadystatechange = function(){
        if (httpRequest.readyState == 4 && httpRequest.status == 200){
            var resultJson = JSON.parse(httpRequest.responseText);
            if(resultJson.result_code == 200){
                location.href = location.origin;
                alert('피드백 감사합니다.');
            }else{
                alert('피드백 전송에 실패하였습니다.');
            }
        }
    };
    httpRequest.open('POST', location.origin + '/api/feedback', true);
    httpRequest.setRequestHeader("Content-type", "application/json");
    httpRequest.send(JSON.stringify(body));
}

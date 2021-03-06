var users;

function getUsers(){
    var xmlHttp = new XMLHttpRequest();
    var searchtext = escape(document.getElementById("cwi").value);

    xmlHttp.onreadystatechange = function() { 
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
            usersFunction(xmlHttp.responseText);
    }

    xmlHttp.open("POST", "/usersget", true); // true for asynchronous
    xmlHttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xmlHttp.send("search="+searchtext);
}

function usersFunction(thetext) {
    var obj = JSON.parse(thetext);
    if(obj == null){
        document.getElementById("cwij").innerHTML="";
        return;
    }
    if(obj.length < 1){
        document.getElementById("cwij").innerHTML="";
        return;
    }

    var thestr = "";
    for(var i= 0; i < obj.length; i++){
        thestr += "<a href=\"javascript:setuser('";
        thestr += obj[i].name + "', " + obj[i].id + ")\">" + obj[i].name + "</a><br />";
    }
    document.getElementById("cwij").innerHTML=thestr;
}

function setuser(thename, theid){
    document.getElementById("cw").textContent=thename;
    document.getElementById("cwname").value=thename;
    document.getElementById("cwnum").value=theid;
}
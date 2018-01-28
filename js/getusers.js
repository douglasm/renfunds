var users;
// window.onload = function(e){ 
//     getUsers(0);
// }

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
        // wname = document.getElementById("origcw").value;
        // val = document.getElementById("origcwnum").value;
        // document.getElementById("cw").textContent=wname;
        // document.getElementById("cwnum").value=val;
        document.getElementById("cwij").textContent="";
        return;
    }
    if(obj.length < 1){
        // wname = document.getElementById("origcw").value;
        // val = document.getElementById("origcwnum").value;
        // document.getElementById("cw").textContent=wname;
        // document.getElementById("cwnum").value=val;
        document.getElementById("cwij").textContent="";
        return;
    }

    var thestr = "";
    for(var i= 0; i < obj.length; i++){
        thestr += "<a href=\"javascript:setuser('";
        thestr += obj[i].name + "', " + obj[i].id + ")\">" + obj[i].name + "</a>" + obj[i].id + "<br />";
    }
    document.getElementById("cwij").innerHTML=thestr;
}

function setuser(thename, theid){
    document.getElementById("cw").textContent=thename;
    document.getElementById("cwname").value=thename;
    document.getElementById("cwnum").value=theid;
}
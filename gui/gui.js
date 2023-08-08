const searchBar = document.getElementById("searchBar");


function stringToUTC(timeString) {
	let time = new Date(timeString);
  
  return time.toUTCString();
}

function clearSearch() {
  if(searchBar.value != "") {
  	searchBar.value = "";
    search();
  }
}

function search() {
    // Declare variables
    var filter, table, tr, td, i, t;
    filter = searchBar.value.toLowerCase();
    table = document.getElementById("buildsListTable");
    tr = table.getElementsByTagName("tr");

    // Loop through all table rows (excluding the header), and hide those who don't match the search query
    for (i = 1; i < tr.length; i++) {
        var filtered = false;
        var tds = tr[i].getElementsByTagName("td");
        for(t = 0; t < tds.length; t++) {
            var td = tds[t];

            if (td) {
              if (td.innerHTML.toLowerCase().indexOf(filter) > -1) {
                filtered = true;
              }
            }     
        }
        if(filtered === true) {
            tr[i].style.display = '';
        }
        else {
            tr[i].style.display = 'none';
        }
    }
}

const buildsTable = document.getElementById("buildsListTableBody");
var builds;

function createBuildsListTable(buildsJSON) {
    builds = JSON.parse(buildsJSON);
    let text = ""

    for (let build in builds) {
        text += "</tr>"

        let tdString = "<td class='buildsListTableCell' onclick='displayDetailsPage(`" + builds[build].buildHash + "`)'>"
        let dateObj = Date.parse(builds[build].time)
        let time = new Date(dateObj).toUTCString();

        text += tdString + time + "</td>"
        text += tdString + builds[build].user + "</td>"
        text += tdString + builds[build].artifact + "</td>"
        text += tdString + builds[build].projectName + "</td>"
        text += tdString + builds[build].gitHash + "</td>"

        text += "</tr>"
    }

    return text
}

const renderBuildsList = async () => {
    let buildsJSON = await getBuildsJSON();
    buildsTable.innerHTML = createBuildsListTable(buildsJSON);
};

function displayHomePage() {
    let detailspage = document.getElementById("detailspage");
    let backBtn = document.getElementById("headerBackBtn");
    let homepage = document.getElementById("homepage");
    
    detailspage.style.display = "none";
    detailspage.style.visibility = "hidden";

    backBtn.style.display = "none";
    backBtn.style.visibility = "hidden";
    
    homepage.style.display = "block";
    homepage.style.visibility = "visible";

    renderBuildsList();
}

async function onPageLoad() {
    // Load and display Builder logo
    let image = await getImage();
    document.getElementById("logo").src = "data:image/png;base64," + image;

    // Display homepage only
    document.getElementById("homepage").style.display = "block";
	document.getElementById("homepage").style.visibility = "visible";
    document.getElementById("detailspage").style.display = "none";
    document.getElementById("detailspage").style.visibility = "hidden";
    document.getElementById("headerBackBtn").style.display = "none";
    document.getElementById("headerBackBtn").style.visibility = "hidden";

    // Render data
    renderBuildsList();
}

// Details page functions

function displayDetailsPage(buildHash) {
    let detailspage = document.getElementById("detailspage");
    let backBtn = document.getElementById("headerBackBtn");
    let homepage = document.getElementById("homepage");
    
    detailspage.style.display = "block";
    detailspage.style.visibility = "visible";

    backBtn.style.display = "inline";
    backBtn.style.visibility = "visible";
    backBtn.innerText = "<";
    
    homepage.style.display = "none";
    homepage.style.visibility = "hidden";

    displayDetailsData(buildHash);
}

async function displayDetailsData(buildHash) {
    let build = builds.find(build => build.buildHash.match(buildHash));

    document.getElementById("projectName").innerHTML = build.projectName;

    let dateObj = Date.parse(build.time)
    let time = new Date(dateObj).toUTCString();
    document.getElementById("timestamp").innerHTML = time;

    // Display metadata
    document.getElementById("projectType").innerHTML = build.projectType;
    document.getElementById("username").innerHTML = build.user;
    document.getElementById("homeDir").innerHTML = build.homeDir;
    document.getElementById("ipAddr").innerHTML = build.ipAddr;
    document.getElementById("gitURL").innerHTML = build.gitURL;
    document.getElementById("gitHash").innerHTML = build.gitHash;
    document.getElementById("branchName").innerHTML = build.branchName;
    document.getElementById("branchHash").innerHTML = build.branchHash;
    
    // Display artifact(s)
    let text = "";
    text += "</tr>"
    text += "<td class='artifact'>" + build.artifact + "</td>";
    text += "</tr>"
    
    document.getElementById("artifactsTableBody").innerHTML = text;
    
    // Display artifact(s) location
    document.getElementById("artifactsLocation").innerHTML = build.artifactLocation;

    // Get logs
    let path = build.artifactLocation;
    path = path.substring(0, path.lastIndexOf('/')) + '/logs/logs.json';
    displayLogs(path);
}

async function displayLogs(path) {
    let logsJSON = await getLogsJSON(path);
	var buildLogs = JSON.parse(logsJSON);
  
    let text = ""
    for (let log in buildLogs) {
        text += "<tr><td>"  	
        
        // Time
        text += "<span class='logTime'>" + new Date(buildLogs[log].timestamp).toLocaleString().replaceAll(',','') + "</span>";
        
        // Type
        if (buildLogs[log].level == 'info'){
            text += "<span class='logTypeInfo'>" + buildLogs[log].level.toUpperCase() + "</span>";
        } else if(buildLogs[log].level == 'warn'){
            text += "<span class='logTypeWarn'>" + buildLogs[log].level.toUpperCase() + "</span>";
        } else {
            text += "<span class='logTypeError'>" + buildLogs[log].level.toUpperCase() + "</span>";
        }
        
        // Caller
        text += "<span class='logCaller'>" + buildLogs[log].caller + ":" + "</span>";
        
        // Message
        text += "<span class='logMessage'>" + buildLogs[log].msg + "</span>";
        
        text += "</td></tr>"
    }
    document.getElementById("logsTable").innerHTML = text;
}
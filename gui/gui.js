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
    var filters, table, tr, td, i, t;
    filters = searchBar.value.toLowerCase().match(/\b(\w+)\b/g);
    table = document.getElementById("buildsListTable");
    tr = table.getElementsByTagName("tr");

    // If filters are deleted repopulate table
    if(searchBar.value == ""){
        renderBuildsList()
    } else {
        // Loop through all table rows (excluding the header), and hide those who don't match the search query
        for (index in filters) {
            for (i = 1; i < tr.length; i++) {
                var filtered = false;
                var tds = tr[i].getElementsByTagName("td");
                for(t = 1; t < tds.length; t++) {
                    var td = tds[t];
        
                    if (td) {
                    if (td.innerHTML.toLowerCase().indexOf(filters[index]) > -1) {
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
    }
}

const buildsTable = document.getElementById("buildsListTableBody");
var builds;

function createBuildsListTable(buildsJSON) {
    builds = JSON.parse(buildsJSON);
    let text = ""

    for (let build in builds.reverse()) {
        text += "<tr>"

        let tdString = "<td class='buildsListTableCell' onclick='displayDetailsPage(`" + builds[build].BuildID + "`)'>"
        let dateObj = Date.parse(builds[build].EndTime)
        let time = new Date(dateObj).toUTCString();

        // Determine what language logo to show
        let image
        switch (builds[build].ProjectType.toLowerCase()) {
            case "c":
                image = "<img src='' alt='C Logo' class='c_logo' width='30' height='30'>"
                break;
            case "csharp":
                image = "<img src='' alt='C# Logo' class='csharp_logo' width='30' height='30'>"
                break;
            case "go":
                image = "<img src='' alt='Go Logo' class='go_logo' width='48' height='30'>"
                break;
            case "java":
                image = "<img src='' alt='Java Logo' class='java_logo' width='30' height='30'>"
                break;
            case "node":
                image = "<img src='' alt='Node Logo' class='node_logo' width='30' height='30'>"
                break;
            case "python":
                image = "<img src='' alt='Python Logo' class='python_logo' width='30' height='30'>"
                break;
            case "ruby":
                image = "<img src='' alt='Ruby Logo' class='ruby_logo' width='30' height='30'>"
                break;
            case "rust":
                image = "<img src='' alt='Rust Logo' class='rust_logo' width='30' height='30'>"
                break;
            default:
                image = ""
        }
        text += tdString + image + "</td>"

        // Fill in table data
        text += tdString + time + "</td>"
        text += tdString + builds[build].UserName + "</td>"
        text += tdString + builds[build].ArtifactName + "</td>"
        text += tdString + builds[build].ProjectName + "</td>"
        text += tdString + builds[build].MasterGitHash + "</td>"

        text += "</tr>"
    }

    return text
}

const renderBuildsList = async () => {
    // Render builds data into table
    let buildsJSON = await getBuildsJSON();
    let tableBodyToDisplay = createBuildsListTable(buildsJSON);
    if (tableBodyToDisplay != "") { 
        // Hide no builds to display div and display table data
        document.getElementById("noBuildsPresentContainer").style.display = "none";
        document.getElementById("noBuildsPresentContainer").style.visibility = "hidden";
        buildsTable.innerHTML = tableBodyToDisplay;
    } else {
        // Show no builds to display div
        document.getElementById("noBuildsPresentContainer").style.display = "block";
        document.getElementById("noBuildsPresentContainer").style.visibility = "visible";
    }

    // Load and display language logos to builds list
    // C
    let c_logo_img = await getCLogoImage();
    const c_logos = document.getElementsByClassName("c_logo");
    for (let i = 0; i < c_logos.length; i++) {
        c_logos[i].src = "data:image/png;base64," + c_logo_img;
    }

    // C#
    let csharp_logo_img = await getCSharpLogoImage();
    const csharp_logos = document.getElementsByClassName("csharp_logo");
    for (let i = 0; i < csharp_logos.length; i++) {
        csharp_logos[i].src = "data:image/png;base64," + csharp_logo_img;
    }

    // Go
    let go_logo_img = await getGoLogoImage();
    const go_logos = document.getElementsByClassName("go_logo");
    for (let i = 0; i < go_logos.length; i++) {
        go_logos[i].src = "data:image/png;base64," + go_logo_img;
    }

    // Java
    let java_logo_img = await getJavaLogoImage();
    const java_logos = document.getElementsByClassName("java_logo");
    for (let i = 0; i < java_logos.length; i++) {
        java_logos[i].src = "data:image/png;base64," + java_logo_img;
    }

    // Node
    let node_logo_img = await getNodeLogoImage();
    const node_logos = document.getElementsByClassName("node_logo");
    for (let i = 0; i < node_logos.length; i++) {
        node_logos[i].src = "data:image/png;base64," + node_logo_img;
    }

    // Python
    let python_logo_img = await getPythonLogoImage();
    const python_logos = document.getElementsByClassName("python_logo");
    for (let i = 0; i < python_logos.length; i++) {
        python_logos[i].src = "data:image/png;base64," + python_logo_img;
    }

    // Ruby
    let ruby_logo_img = await getRubyLogoImage();
    const ruby_logos = document.getElementsByClassName("ruby_logo");
    for (let i = 0; i < ruby_logos.length; i++) {
        ruby_logos[i].src = "data:image/png;base64," + ruby_logo_img;
    }

    //Rust
    let rust_logo_img = await getRustLogoImage();
    const rust_logos = document.getElementsByClassName("rust_logo");
    for (let i = 0; i < rust_logos.length; i++) {
        rust_logos[i].src = "data:image/png;base64," + rust_logo_img;
    }
    
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

    // Render builds data
    renderBuildsList();
}

// Details page functions

function displayDetailsPage(buildID) {
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

    displayDetailsData(buildID);
}

async function displayDetailsData(buildID) {
    let build = builds.find(build => build.BuildID.match(buildID));

    document.getElementById("projectName").innerHTML = build.ProjectName;

    let dateObj = Date.parse(build.EndTime)
    let time = new Date(dateObj).toUTCString();
    document.getElementById("timestamp").innerHTML = time;

    // Display metadata
    document.getElementById("projectType").innerHTML = build.ProjectType;
    document.getElementById("username").innerHTML = build.UserName;
    document.getElementById("homeDir").innerHTML = build.HomeDir;
    document.getElementById("ipAddr").innerHTML = build.IP;
    document.getElementById("gitURL").innerHTML = build.GitURL;
    document.getElementById("gitHash").innerHTML = build.MasterGitHash;
    document.getElementById("branchName").innerHTML = build.BranchName;
    document.getElementById("buildID").innerHTML = build.BuildID;
    
    // Display artifact(s)
    let text = "";

    let artifactArray = build.ArtifactName.split(",")

    for (artifact in artifactArray) {
        text += "</tr>"
        text += "<td class='artifact'>" + artifactArray[artifact] + "</td>";
        text += "</tr>"
    }
    
    document.getElementById("artifactsTableBody").innerHTML = text;
    
    // Display artifact(s) location
    document.getElementById("artifactsLocation").innerHTML = build.ArtifactLocation;

    // Get logs
    let path = build.LogsLocation;
    //path = path.substring(0, path.lastIndexOf('/')) + '/logs/logs.json';
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
        //text += "<span class='logCaller'>" + buildLogs[log].caller + ":" + "</span>";
        
        // Message
        text += "<span class='logMessage'>" + buildLogs[log].msg + "</span>";
        
        text += "</td></tr>"
    }
    document.getElementById("logsTable").innerHTML = text;
}
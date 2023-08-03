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

const render = async () => {
    buildsTable.innerHTML = await jsonToHTML();
};

async function onPageLoad() {
    // Load and display Builder logo
    let image = await getImage();
    document.getElementById("logo").src = "data:image/png;base64," + image;

    render();
}
function process() {
	console.log("The button was clicked.");
	var xhr = new XMLHttpRequest();
	xhr.open("GET","/getnumber", true);
	xhr.onreadystatechange = function() {
		console.log(xhr.readyState);
		if (xhr.readyState == 1) {
			document.getElementById("changeme").innerHTML = "loading...";
			console.log("Ajax loading.");
		}
		if (xhr.readyState == 4 && xhr.status == 200) {
			console.log("Ajax ready.");
			document.getElementById("changeme").innerHTML = xhr.responseText;
		}

		if (xhr.readyState == 4 && xhr.status != 200) {
			document.getElementById("changeme").innerHTML = "No lucky number available right now. Please try again later.";
		}
	};
	xhr.send();
	document.getElementById("changeme").innerHTML = "loading...";
}

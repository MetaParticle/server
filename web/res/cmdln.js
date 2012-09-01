"use strict";

$(document).ready(function() {
	$("#form").on("submit", function(event) {
		console.log($(this).text());
		writeToOutput($("#in").value, "blue");
		return false;
	});
});

function writeToOutput(string, color) {
	//$("#output").append(string);
}

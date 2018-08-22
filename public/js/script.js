// request address data from server
var addresses = $.ajax({
    url: "http://127.0.0.1:8000/api/v1/addresses",
    // url: "http://127.0.0.1:8000/api/v1/addresses?bbox=-79,35,-78,36",
    dataType: "json",
    error: function(xhr) {
        alert(xhr.statusText)
    }
})
$.when(addresses).done(function() {
    // create map
    var map = L.map('map')
        .setView([35.787743, -78.644257], 7);

    var basemap = L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png', {
        attribution: 'Map data &copy; <a href="https://openstreetmap.org">OpenStreetMap</a> contributors',
        maxZoom: 15
    }).addTo(map);

    // add addresses to map
    L.heatLayer(addresses.responseJSON).addTo(map)
});

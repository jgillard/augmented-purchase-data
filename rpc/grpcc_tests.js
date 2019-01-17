client.getStatus({}, function(err, data) {
    console.log("Ran getStatus({})")
    console.log("", data);
    if (data.status !== "OK") {
        console.log("expected {status: 'OK'}")
    } else {
        console.log("PASS")
    }
});

client.ListCategories({}, function(err, data) {
    console.log("Ran ListCategories({})")
    console.log("", data);
});

client.GetCategory({ categoryID: "1" }, function(err, data) {
    console.log("Ran GetCategory({categoryID:'1'})")
    console.log("", data);
});

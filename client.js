(async () => {
    function addFileDownload(filename, id, fileURL) {
        let newFilename = ""

        if(filename == null) {
            console.log("Skipping image you don't have access to...")
        } else {
            let originalFilename = filename.split(".")
            let fileExtension = originalFilename.pop() // retrieves only the extension, text after last dot
            newFilename = `${originalFilename.join(".")}-${id}.${fileExtension}` // frankensteins together a filename with -<id> appended
        }

        downloads.push([
            fileURL,
            newFilename
        ])
    }

    let basePostURL = "https://www.patreon.com/api/posts?"

    let campaignID = Number(window.patreon.bootstrap.creator.data.id)

    console.log("Sending Patreon information to patreon-dl...")

    await fetch("http://localhost:9849/user", {
        method: 'POST',
        body: JSON.stringify({
            id: campaignID,
            name: window.patreon.bootstrap.creator.data.attributes.name
        })
    });

    const initialQueryParams = new URLSearchParams({
        "include": "images,media",
        "fields[post]": "post_metadata",
        "fields[media]": "id,image_urls,download_url,metadata,file_name",
        "filter[campaign_id]": campaignID,
        "filter[contains_exclusive_posts]": true,
        "sort": "-published_at",
        "json-api-version": "1.0"
    })

    let downloads = [];

    let initalPostRequest = await fetch(basePostURL + initialQueryParams.toString())
    let parsedInital = await initalPostRequest.json()

    let initialLength = 0;
    if ("included" in parsedInital) {
        initialLength = parsedInital.included.length
    }

    for (let i = 0; i < initialLength; i++) {
        let originalFilename = parsedInital.included[i].attributes.file_name.split(".")
        let fileExtension = originalFilename.pop() // retrieves only the extension, text after last dot
        let newFilename = `${originalFilename.join(".")}-${parsedInital.included[i].id}.${fileExtension}` // frankensteins together a filename with -<id> appended

        addFileDownload(parsedInital.included[i].attributes.file_name, parsedInital.included[i].id, parsedInital.included[i].attributes.download_url);
    }

    console.log("Collected " + downloads.length + " posts...")

    let nextURL = ""
    if ("links" in parsedInital) {
        nextURL = parsedInital.links.next;
    }

    while (nextURL !== "") {
        let recursivePostRequest = await fetch(nextURL)
        let parsedPosts = await recursivePostRequest.json()

        let includedLength = 0;
        if ("included" in parsedPosts) {
            includedLength = parsedPosts.included.length
        }

        for (let i = 0; i < includedLength; i++) {
            addFileDownload(parsedPosts.included[i].attributes.file_name, parsedPosts.included[i].id, parsedPosts.included[i].attributes.download_url);
        }

        if ("links" in parsedPosts) {
            nextURL = parsedPosts.links.next;
        } else {
            nextURL = ""
        }
        console.log("Collected " + downloads.length + " posts...")
    }

    console.log("Sending " + downloads.length + " image links to patreon-dl...")
    await fetch("http://localhost:9849/download", {
        method: 'POST',
        body: JSON.stringify(downloads)
    });

    await fetch("http://localhost:9849/done");
    console.log("patreon-dl is starting download...")
})();
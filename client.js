(async () => {
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
        downloads.push([
            parsedInital.included[i].attributes.download_url,
            parsedInital.included[i].attributes.file_name
        ])
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
            downloads.push([
                parsedPosts.included[i].attributes.download_url,
                parsedPosts.included[i].attributes.file_name
            ])
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
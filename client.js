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

    console.log("%cpatreon-dl Downloader - Starting...", "font-size:large")

    if(!document.location.pathname.endsWith("/posts")) {
        console.error("Invalid URL. Are you sure you are on the /posts page of the creator?\nStopping...")
        return
    }

    const basePostURL = "https://www.patreon.com/api/posts?"

    const campaignID = Number(window.patreon.bootstrap.creator.data.id) // this is an internal ID for the campaign, NOT the username

    console.log("Sending Patreon information to patreon-dl...")

    await fetch("http://localhost:9849/user", {
        method: 'POST',
        body: JSON.stringify({
            id: campaignID,
            name: window.patreon.bootstrap.creator.data.attributes.name
        })
    });

    // this filters the posts to only get posts that have some form of media + a download url for that media attached
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

    const initalPostRequest = await fetch(basePostURL + initialQueryParams.toString())
    const parsedInital = await initalPostRequest.json()

    let initialLength = 0;
    if ("included" in parsedInital) {
        initialLength = parsedInital.included.length
    }

    for (let i = 0; i < initialLength; i++) {
        if(parsedInital.included[i].attributes.file_name === null) {
            continue
        }

        const originalFilename = parsedInital.included[i].attributes.file_name.split(".")
        const fileExtension = originalFilename.pop() // retrieves only the extension, text after last dot
        const newFilename = `${originalFilename.join(".")}-${parsedInital.included[i].id}.${fileExtension}` // frankensteins together a filename with -<id> appended

        addFileDownload(parsedInital.included[i].attributes.file_name, parsedInital.included[i].id, parsedInital.included[i].attributes.download_url);
    }

    console.log(`Collected ${downloads.length} posts...`)

    let nextURL = ""
    if ("links" in parsedInital) {
        nextURL = parsedInital.links.next;
    }

    while (nextURL !== "") {
        await new Promise(r => setTimeout(r, 1000)); // wait for 1s to prevent #11

        const recursivePostRequest = await fetch(nextURL)
        const parsedPosts = await recursivePostRequest.json()

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
        console.log(`%cCollected %c${downloads.length} %cposts...`, "color:white", "color:green", "color:white");
    }

    console.log("%cEnd of posts reached (total: %d)!", "color:green", downloads.length);

    console.log(`Sending ${downloads.length} image links to patreon-dl...`)
    await fetch("http://localhost:9849/download", {
        method: 'POST',
        body: JSON.stringify(downloads)
    });

    console.log("patreon-dl is starting download...")
})();
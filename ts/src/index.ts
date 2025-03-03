import Uppy from "@uppy/core";
import Tus from "@uppy/tus"
import Dashboard from "@uppy/dashboard"

const uppyContainer = document.getElementById("uppy")

if(uppyContainer != null && !uppyContainer.innerHTML) {
    const uppy = new Uppy()
    uppy.use(Dashboard, { target: '#uppy', inline: true, fileManagerSelectionType: 'both' }).use(Tus, {endpoint: 'http://localhost:8081/uploads/'})
}

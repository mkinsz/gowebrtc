let pc = new RTCPeerConnection()

pc.onnegotiationneeded = handleNegotiationNeededEvent;
pc.oniceconnectionstatechange = e => console.log('ICE State: ', pc.iceConnectionState)

pc.ontrack = function (event) {
    var el = document.createElement(event.track.kind)
    el.srcObject = event.streams[0]
    el.muted = true
    el.autoplay = true
    el.controls = true
    el.width = 600

    document.getElementById('remoteVideos').appendChild(el)
}

pc.addTransceiver('video', {
    'direction': 'sendrecv'
})

async function handleNegotiationNeededEvent() {
    pc.createOffer()
        .then(offer => {
            pc.setLocalDescription(offer)

            return fetch(`/signal`, {
                method: 'post',
                headers: {
                    'Accept': 'application/json, text/plain, */*',
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(offer)
            })
        })
        .then(res => res.json())
        .then(res => pc.setRemoteDescription(res))
        .catch(alert)
}

$(document).ready(function () {
    fetch(`/codec`).then(resp => resp.json()).then(
        data => {
        try {
            console.log("Codec: ", data)
        } catch (e) {
            console.log('Fetch Codec err: ', e)
        } finally {
            pc.addTransceiver('video', {
                'direction': 'sendrecv'
            });
        }
    })
})
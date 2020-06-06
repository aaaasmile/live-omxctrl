
const handleError = (error, that) => {
	console.error(error);
	if (error.bodyText !== '') {
		that.$store.commit('msgText', `${error.statusText}: ${error.bodyText}`)
	} else {
		that.$store.commit('msgText', 'Error: empty response')
	}
}

export default {
	TogglePowerState(that, req) {
		console.log('Request is ', req)
		that.$http.post("TogglePowerState", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			if (req.power == 'off'){
				that.poweron = false
			}else if(req.power == 'on'){
				that.poweron = true
			}
		}, error => {
			handleError(error, that)
		});
	},
	ChangeVolume(that, req) {
		console.log('Request is ', req)
		that.$http.post("ChangeVolume", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			if (req.volume == 'mute'){
				that.muted = true
			}else if(req.volume == 'unmute'){
				that.muted = false
			}
		}, error => {
			handleError(error, that)
		});
	},
	PlayURI(that, req) {
		console.log('Request is ', req)
		that.$http.post("PlayURI", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.curruri = req.URI
			that.loadingMeta = false
		}, error => {
			that.loadingMeta = false
			handleError(error, that)
		});
	},
	Resume(that, req) {
		console.log('Request is ', req)
		that.$http.post("PlayURI", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.loadingMeta = false
		}, error => {
			that.loadingMeta = false
			handleError(error, that)
		});
	},
	Pause(that, req) {
		console.log('Request is ', req)
		that.$http.post("Pause", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.loadingMeta = false
		}, error => {
			that.loadingMeta = false
			handleError(error, that)
		});
	},
	GetPlayerState(that, req) {
		console.log('Request is ', req)
		that.$http.post("GetPlayerState", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.$store.commit('playerstate', result.data)
		}, error => {
			handleError(error, that)
		});
	},
}
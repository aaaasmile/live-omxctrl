
const handleError = (error, that) => {
	console.error(error);
	that.loadingMeta = false
	if (error.bodyText !== '') {
		that.$store.commit('msgText', `${error.statusText}: ${error.bodyText}`)
	} else {
		that.$store.commit('msgText', 'Error: empty response')
	}
}

export default {
	SetPowerState(that, req) {
		console.log('Request is ', req)
		that.$http.post("SetPowerState", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.$store.commit('playerstate', result.data)
			that.loadingMeta = false
		}, error => {
			handleError(error, that)
		});
	},
	ChangeVolume(that, req) {
		console.log('Request is ', req)
		that.$http.post("ChangeVolume", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.$store.commit('playerstate', result.data)
		}, error => {
			handleError(error, that)
		});
	},
	PlayURI(that, req) {
		console.log('Request is ', req)
		that.$http.post("PlayURI", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.$store.commit('playerstate', result.data)
		}, error => {
			handleError(error, that)
		});
	},
	Resume(that, req) {
		console.log('Resume Request is ', req)
		that.$http.post("PlayURI", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.$store.commit('playerstate', result.data)
		}, error => {
			handleError(error, that)
		});
	},
	Pause(that, req) {
		console.log('Pause Request is ', req)
		that.$http.post("Pause", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.$store.commit('playerstate', result.data)
		}, error => {
			handleError(error, that)
		});
	},
	GetPlayerState(that, req) {
		console.log('Request is ', req)
		that.$http.post("GetPlayerState", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.$store.commit('playerstate', result.data)
			that.loadingMeta = false
		}, error => {
			handleError(error, that)
		});
	},
	NextTitle(that, req) {
		console.log('Request is ', req)
		that.$http.post("NextTitle", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.$store.commit('playerstate', result.data)
			that.loadingMeta = false
		}, error => {
			handleError(error, that)
		});
	}
}
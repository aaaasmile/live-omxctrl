
const handleError = (error, that) => {
	console.error(error);
	if (error.bodyText !== '') {
		that.$store.commit('msgText', `${error.statusText}: ${error.bodyText}`)
	} else {
		that.$store.commit('msgText', 'Error: empty response')
	}
}

export default {
	ChangeVolume(that, req) {
		console.log('Request is ', req)
		that.$http.post("ChangeVolume", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.$store.commit('msgText', result.data.Status)
		}, error => {
			handleError(error, that)
		});
	},
	PlayURI(that, req) {
		console.log('Request is ', req)
		that.$http.post("PlayURI", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			//that.$store.commit('msgText', result.data.Status)
			that.loadingMeta = false
		}, error => {
			that.loadingMeta = false
			handleError(error, that)
		});
	},
	PauseURI(that, req) {
		console.log('Request is ', req)
		that.$http.post("PauseURI", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			//that.$store.commit('msgText', result.data.Status)
			that.loadingMeta = false
		}, error => {
			that.loadingMeta = false
			handleError(error, that)
		});
	},
	Stop(that, req) {
		console.log('Request is ', req)
		that.$http.post("Stop", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.$store.commit('msgText', result.data.Status)
		}, error => {
			handleError(error, that)
		});
	},
}
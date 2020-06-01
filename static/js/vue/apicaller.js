
const handleError = (error, that) => {
	console.error(error);
	if (error.bodyText !== '') {
		that.$store.commit('msgText', `${error.statusText}: ${error.bodyText}`)
	} else {
		that.$store.commit('msgText', 'Error: empty response')
	}
}

export default {
	ChangeVolume(that, req, scope) {
		req.Scope = scope
		console.log('Request is ', req)
		that.$http.post("ChangeVolume", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.$store.commit('msgText', result.data.Status)
		}, error => {
			handleError(error, that)
		});
	},
	Play(that, req, scope) {
		req.Scope = scope
		console.log('Request is ', req)
		that.$http.post("Play", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.$store.commit('msgText', result.data.Status)
			that.loadingMeta = false
		}, error => {
			that.loadingMeta = false
			handleError(error, that)
		});
	},
	Stop(that, req, scope) {
		req.Scope = scope
		console.log('Request is ', req)
		that.$http.post("Stop", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call terminated ', result.data)
			that.$store.commit('msgText', result.data.Status)
		}, error => {
			handleError(error, that)
		});
	},
}
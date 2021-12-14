
const handleError = (error, that) => {
	console.error(error);
	that.loadingMeta = false
	if (error.bodyText !== '') {
		that.$store.commit('msgText', `${error.statusText}: ${error.bodyText}`)
	} else {
		that.$store.commit('msgText', 'Error: empty response')
	}
}

const closeAllDialogs = (that) => {
	that.osloading = false
	that.dialogShutdown = false
	that.dialogReboot = false
	that.dialogRestartService = false
}

export default {
	SetPowerState(that, req) {
		console.log('Request is ', req)
		that.$http.post("SetPowerState", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.$store.commit('playerstate', result.data)
			that.loadingMeta = false
		}, error => {
			handleError(error, that)
		});
	},
	ChangeVolume(that, req) {
		console.log('Request is ', req)
		that.$http.post("ChangeVolume", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.$store.commit('playerstate', result.data)
		}, error => {
			handleError(error, that)
		});
	},
	Resume(that, req) {
		console.log('Resume Request is ', req)
		that.$http.post("Resume", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.$store.commit('playerstate', result.data)
		}, error => {
			handleError(error, that)
		});
	},
	Pause(that, req) {
		console.log('Pause Request is ', req)
		that.$http.post("Pause", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.$store.commit('playerstate', result.data)
		}, error => {
			handleError(error, that)
		});
	},
	GetPlayerState(that, req) {
		console.log('Request is ', req)
		that.$http.post("GetPlayerState", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.$store.commit('playerstate', result.data)
			that.loadingMeta = false
		}, error => {
			handleError(error, that)
		});
	},
	NextTitle(that, req) {
		console.log('Request is ', req)
		that.$http.post("NextTitle", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.$store.commit('playerstate', result.data)
			that.loadingMeta = false
		}, error => {
			handleError(error, that)
		});
	},
	PreviousTitle(that, req) {
		console.log('Request is ', req)
		that.$http.post("PreviousTitle", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.$store.commit('playerstate', result.data)
			that.loadingMeta = false
		}, error => {
			handleError(error, that)
		});
	},
	PlayUri(that, req) {
		console.log('Request is ', req)
		that.$http.post("PlayUri", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.$store.commit('playerstate', result.data)
			that.loadingyoutube = false
		}, error => {
			that.loadingyoutube = false
			handleError(error, that)
		});
	},
	OSRequest(that, req) {
		console.log('OS Request is ', req)
		that.osloading = true
		that.$http.post("OSRequest", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			closeAllDialogs(that)	
			that.$store.commit('msgText', result.data.msg)
		}, error => {
			closeAllDialogs(that)	
			handleError(error, that)
		});
	},
	FetchHistory(that, req) {
		console.log('FetchHistory request is ', req)
		that.hisloading = true
		that.$http.post("FetchHistory", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.hisloading = false
			that.$store.commit('historyfetch', result.data)
		}, error => {
			that.hisloading = false	
			handleError(error, that)
		});
	},
	FetchVideo(that, req) {
		req.name = 'FetchVideo'
		console.log('FetchVideo request is ', req)
		that.videoloading = true
		that.$http.post("HandleVideo", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.videoloading = false
			that.$store.commit('videofetch', result.data)
		}, error => {
			that.videoloading = false	
			handleError(error, that)
		});
	},
	ScanVideo(that, req) {
		req.name = 'ScanVideo'
		console.log('ScanVideo request is ', req)
		that.videoloading = true
		that.dialogScan = false
		that.$http.post("HandleVideo", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.videoloading = false
			that.$store.commit('videofetch', result.data)
		}, error => {
			that.videoloading = false	
			handleError(error, that)
		});
	},
	ScanMusic(that, req, fn) {
		req.name = 'ScanMusic'
		console.log('ScanMusic request is ', req)
		that.loadingData = true
		that.dialogScan = false
		that.$http.post("HandleMusic", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.loadingData = false
			that.$store.commit('musicfetch', result.data)
			if(fn){
				fn()
			}
		}, error => {
			that.loadingData = false	
			handleError(error, that)
			if(fn){
				fn()
			}
		});
	},
	FetchMusic(that, req, fn) {
		req.name = 'FetchMusic'
		console.log('FetchMusic request is ', req)
		that.loadingData = true
		that.$http.post("HandleMusic", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.loadingData = false
			that.$store.commit('musicfetch', result.data)
			that.page = 1
			that.search = ''
			if(fn){
				fn()
			}
		}, error => {
			that.loadingData = false	
			handleError(error, that)
			if(fn){
				fn()
			}
		});
	},
	FetchRadio(that, req) {
		req.name = 'FetchRadio'
		console.log('FetchRadio request is ', req)
		that.radioloading = true
		that.$http.post("HandleRadio", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.radioloading = false
			that.$store.commit('radiofetch', result.data)
		}, error => {
			that.radioloading = false	
			handleError(error, that)
		});
	},
	HandleRadio(that, req, fnCompl) {
		console.log('HandleRadio request is ', req)
		that.radioloading = true
		that.$http.post("HandleRadio", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.radioloading = false
			if(fnCompl){
				fnCompl(true,result)
			}
		}, error => {
			that.radioloading = false	
			handleError(error, that)
			if(fnCompl){
				fnCompl(false)
			}
		});
	}
}
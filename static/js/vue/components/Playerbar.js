import API from '../apicaller.js'

export default {
	components: {},
	data() {
		return {
			playing: false,
			loadingMeta: false
		}
	},
	computed: {
		...Vuex.mapState({

		}),

	},
	methods: {
		togglePlayURI() {
			if (this.playing) {
				this.playing = !this.playing
				console.log('Pause URI')
				API.PauseURI(this, req)
			} else {
				console.log('Play URI')
				let req = {URI: 'http://stream.srg-ssr.ch/m/rsc_de/aacp_96'}
				API.PlayURI(this, req)
			}
		}
	},
	template: `
  <v-toolbar flat>
    <v-tooltip bottom>
      <template v-slot:activator="{ on }">
        <v-btn icon  v-on="on">
          <v-icon>mdi-skip-previous</v-icon>
        </v-btn>
      </template>
      <span>Previous</span>
    </v-tooltip>

    <v-tooltip bottom>
      <template v-slot:activator="{ on }">
        <v-btn icon  v-on="on" @click="togglePlayURI">
          <v-icon>{{ playing ? 'mdi-pause' : 'mdi-play' }}</v-icon>
        </v-btn>
      </template>
      <span>{{ playing ? 'Pause' : 'Play current'}}</span>
    </v-tooltip>

    <v-tooltip bottom>
      <template v-slot:activator="{ on }">
        <v-btn icon  v-on="on">
          <v-icon>mdi-skip-next</v-icon>
        </v-btn>
      </template>
      <span>Next</span>
    </v-tooltip>

    <v-tooltip bottom>
      <template v-slot:activator="{ on }">
        <v-btn icon  v-on="on">
          <v-icon>mdi-shuffle</v-icon>
        </v-btn>
      </template>
      <span>Shuffle</span>
    </v-tooltip>

    <v-tooltip bottom>
      <template v-slot:activator="{ on }">
        <v-btn icon  v-on="on">
          <v-icon>mdi-repeat</v-icon>
        </v-btn>
      </template>
      <span>Repeat</span>
    </v-tooltip>
  </v-toolbar>`
}

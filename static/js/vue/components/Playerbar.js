export default {
    components: {  },
    data() {
      return {
        playing: false,
      }
    },
    computed: {
      ...Vuex.mapState({
        
      }),
     
    },
    methods: {
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
        <v-btn icon  v-on="on" @click="playing = !playing">
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
  </v-toolbar>
`
}
  
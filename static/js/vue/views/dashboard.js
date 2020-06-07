export default {
    data() {
        return {
        }
    },
    computed: {
      ...Vuex.mapState({
        PlayingURI: state => {
          return state.ps.uri
        },
      })
    },
    methods: {
    },
    template: `
  <v-card color="grey lighten-4" flat tile>
    <v-toolbar flat dense>
      <v-toolbar-title class="subheading grey--text">Dashboard for Omx Player</v-toolbar-title>
    </v-toolbar>
    <v-col class="mb-5" cols="12">
      <v-row justify="center">
        <v-col class="mb-5" cols="12">
          <v-card>
            <v-card-title>Current Media</v-card-title>
            <div class="mx-4">
              <div class="subtitle-2 text--secondary">URI: {{PlayingURI}}</div>
            </div>
          </v-card>
        </v-col>
      </v-row>
    </v-col>
  </v-card>`
}
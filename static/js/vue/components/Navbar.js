import routes from '../routes.js'
import Toast from './toast.js'

export default {
  components: { Toast },
  data() {
    return {
      drawer: false,
      links: routes,
      AppTitle: "E-Mailer",
    }
  },
  computed: {
    ...Vuex.mapState({
      
    }),
   
  },
  methods: {
    gotoAllow() {
      console.log("Request auth")
    }
  },
  template: `
  <nav>
    <v-app-bar dense flat>
      <v-btn text color="grey">
        <v-icon>menu</v-icon>
      </v-btn>
      <v-toolbar-title class="text-uppercase grey--text">
        <span class="font-weight-light">Live</span>
        <span>{{AppTitle}}</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <!-- dropdown menu -->
      <v-menu offset-y>
        <template v-slot:activator="{ on }">
          <v-btn v-on="on" text color="grey">
            <v-icon left>expand_more</v-icon>
            <span>Menu</span>
          </v-btn>
        </template>
        <v-list>
          <v-list-item v-for="link in links" :key="link.text" router :to="link.path">
            <v-list-item-title>{{ link.text }}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-app-bar>
    <Toast></Toast>
  </nav>
`
}
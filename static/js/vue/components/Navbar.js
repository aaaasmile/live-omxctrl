import routes from '../routes.js'
import Toast from './toast.js'

export default {
  components: { Toast },
  data() {
    return {
      links: routes,
      AppTitle: "Omx Control",
      drawer: false,
    }
  },
  computed: {
    ...Vuex.mapState({
      
    }),
   
  },
  methods: {
  },
  template: `
  <nav>
    <v-app-bar dense flat>
      <v-btn text color="grey" @click="drawer = !drawer">
        <v-icon>mdi-menu</v-icon>
      </v-btn>
      <v-toolbar-title class="text-uppercase grey--text">
        <span class="font-weight-light">Live</span>
        <span>{{AppTitle}}</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
    </v-app-bar>
    <Toast></Toast>
    <v-navigation-drawer app v-model="drawer">
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="title">{{AppTitle}}</v-list-item-title>
          <v-list-item-subtitle>player</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>

      <v-divider></v-divider>

      <v-list dense nav>
        <v-list-item v-for="item in links" :key="item.title" link>
          <v-list-item-icon>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
  </nav>`
}
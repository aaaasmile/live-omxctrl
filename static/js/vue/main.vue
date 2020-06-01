<template>
  <v-app class="grey lighten-4">
    <Navbar />
    <v-content class="mx-4 mb-4">
      <router-view></router-view>
    </v-content>
    <v-footer absolute>
      <v-container>
        <v-row>
          <v-col class="d-flex text-center caption" cols="12">
            <Playerbar />
          </v-col>
        </v-row>
        <v-row>
          <v-col class="d-flex text-center caption" cols="12">
            <div>
              {{ new Date().getFullYear() }} —
              <span>Buildnr: {{Buildnr}}</span>
            </div>
          </v-col>
        </v-row>
      </v-container>
    </v-footer>
  </v-app>
</template>
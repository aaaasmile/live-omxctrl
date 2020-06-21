<template>
  <v-card color="grey lighten-4" flat tile>
    <v-toolbar flat dense>
      <v-toolbar-title class="subheading grey--text">Dashboard for Omx Player</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-btn icon @click="playYoutube" :loading="loadingyoutube" v-on="on">
            <v-icon>airplay</v-icon>
          </v-btn>
        </template>
        <span>Play youtube url</span>
      </v-tooltip>
    </v-toolbar>
    <v-col class="mb-1" cols="12">
      <v-row justify="center">
        <v-col cols="12">
          <v-text-field v-model="uriToPlay" label="Select an URL"></v-text-field>
        </v-col>
        <v-row>
          <v-col class="mb-1" cols="12">
            <v-card>
              <v-card-title>Current Media</v-card-title>
              <div class="mx-4">
                <div class="subtitle-2 text--secondary">URI: {{PlayingURI}}</div>
              </div>
            </v-card>
          </v-col>
        </v-row>
      </v-row>
      <v-row>
        <v-col cols="12">
          <Playerbar class="mr-1 ml-1"/>
        </v-col>
      </v-row>
    </v-col>
  </v-card>
</template>
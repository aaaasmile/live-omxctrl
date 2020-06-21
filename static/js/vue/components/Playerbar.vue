<template>
  <v-skeleton-loader
    :loading="loadingMeta"
    :transition="transition"
    height="94"
    type="list-item-two-line"
  >
    <v-card flat tile>
      <v-row>
        <v-col class="md-6" cols="12">
          <v-toolbar flat>
            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-btn icon v-on="on" @click="previousTitle">
                  <v-icon>mdi-skip-previous</v-icon>
                </v-btn>
              </template>
              <span>Previous</span>
            </v-tooltip>

            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-btn icon v-on="on" @click="togglePlayResume">
                  <v-icon>{{ Playing ? 'mdi-pause' : 'mdi-play' }}</v-icon>
                </v-btn>
              </template>
              <span>{{ Playing ? 'Pause' : 'Play current'}}</span>
            </v-tooltip>

            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-btn icon v-on="on" @click="nextTitle">
                  <v-icon>mdi-skip-next</v-icon>
                </v-btn>
              </template>
              <span>Next</span>
            </v-tooltip>

            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-btn icon v-on="on">
                  <v-icon>mdi-shuffle</v-icon>
                </v-btn>
              </template>
              <span>Shuffle</span>
            </v-tooltip>

            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-btn icon v-on="on">
                  <v-icon>mdi-repeat</v-icon>
                </v-btn>
              </template>
              <span>Repeat</span>
            </v-tooltip>
          </v-toolbar>
        </v-col>
        <v-col class="md-6" cols="12">
          <v-toolbar flat>
            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-btn icon v-on="on" @click="toggleMute">
                  <v-icon>{{ Muted ? 'volume_mute' : 'volume_off' }}</v-icon>
                </v-btn>
              </template>
              <span>{{ Muted ? 'Unmute' : 'Mute'}}</span>
            </v-tooltip>
            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-btn icon v-on="on" @click="VolumeDown">
                  <v-icon>volume_down</v-icon>
                </v-btn>
              </template>
              <span>Volume down</span>
            </v-tooltip>
            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-btn icon v-on="on" @click="VolumeUp">
                  <v-icon>volume_up</v-icon>
                </v-btn>
              </template>
              <span>Volume up</span>
            </v-tooltip>
            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-btn icon v-on="on" @click="syncStatus">
                  <v-icon>mdi-sync</v-icon>
                </v-btn>
              </template>
              <span>Synchronize status</span>
            </v-tooltip>
            <v-btn icon @click="togglePower" :color="ColorPower">
              <v-icon>power_settings_new</v-icon>
            </v-btn>
          </v-toolbar>
        </v-col>
      </v-row>
    </v-card>
  </v-skeleton-loader>
</template>
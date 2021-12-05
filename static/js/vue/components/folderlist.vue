<template>
  <v-card>
    <v-card-title>
      <v-text-field v-model="search" append-icon="search" label="Search" single-line hide-details></v-text-field>
    </v-card-title>
    <v-data-table
      v-model="musicSelected"
      :headers="headers"
      :items="musicdata"
      :loading="loadingData"
      :items-per-page="itemsPerPage"
      item-key="id"
      show-select
      class="elevation-1"
      :search="search"
      :footer-props="{
      showFirstLastPage: true,
      firstIcon: 'mdi-arrow-collapse-left',
      lastIcon: 'mdi-arrow-collapse-right',
      prevIcon: 'mdi-minus',
      nextIcon: 'mdi-plus'
    }"
    >
      <template v-slot:item.actions="{ item }">
        <v-icon small class="mr-2" @click="playOrfetchSubFolder(item)" >{{item.icon_action}}</v-icon>
      </template>
      <template v-slot:item.type="{ item }">
        <v-chip
          :color="getColorType(item.fileorfolder)"
          dark
        >{{ item.type }}</v-chip>
      </template>
    </v-data-table>
  </v-card>
</template>
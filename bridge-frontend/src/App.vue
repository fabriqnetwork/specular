<script setup lang="ts">
import { RouterLink, RouterView } from "vue-router";
import Header from "./components/Header.vue";
import Footer from "./components/Footer.vue";
import FAQ from "./components/FAQ.vue";
</script>

<template>
  <v-app full-height>
    <Header />
    <v-main>
      <v-container fluid class="fill-height">
        <v-row align="center" justify="center">
          <v-col cols="6" style="max-width: 600px">
            <RouterView />
            <FAQ />
          </v-col>
        </v-row>
      </v-container>
    </v-main>
    <Footer/>
  </v-app>
</template>

<style scoped></style>

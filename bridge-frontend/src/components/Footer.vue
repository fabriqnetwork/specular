<template>
  <v-footer app :class="footer.style">
    <v-row align="center" justify="center" class="footer-row">
      <v-col cols="4" v-for="link in footer.links" :key="link.title" class="text-center">
        <div>
          <h4 class="footer-heading">{{ link.title }}</h4>
          <ul class="footer-list">
            <li v-for="item in link.items" :key="item.label">
              <a :href="item.to || item.href" class="footer-link" target="_blank">{{ item.label }}</a>
            </li>
          </ul>
        </div>
      </v-col>
    </v-row>
    <v-col cols="1" class="text-center">
      <span class="footer-caption">{{ footer.caption }}</span>
    </v-col>
  </v-footer>
</template>

<script>
export default {
  name: 'Footer',
  data() {
    return {
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Getting Started',
                to: 'https://specular.network/docs/getting-started',
              },
              {
                label: 'Community',
                to: 'https://specular.network/docs/community/',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'Twitter',
                href: 'https://twitter.com/specularl2',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'Blog',
                href: 'https://medium.com/@SpecularL2',
              },
              {
                label: 'GitHub',
                href: 'https://github.com/specularl2',
              },
            ],
          },
        ],
        caption: `Copyright © 2023 Specular.`,
      },
    };
  },
};
</script>

<style scoped>
.dark {
  background-color: #333;
  color: #fff;
}

.footer-heading {
  padding: 5px;
  color: #fff;
  margin: 0 auto;
}

.footer-link {
  color: #ccc;
  text-decoration: none;
}

.footer-list {
  list-style-type: none;
  padding-left: 0;
  margin: 0 auto;
}

.footer-row {
  padding-top: 5px;
  padding-bottom: 20px;
}

.footer-caption {
  font-size: 12px;
  opacity: 0.7;
  color: #fff;
  margin-top: 10px; /* Adjust the space as needed */
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
}
</style>

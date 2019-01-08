import Vue from 'vue';
import Vuex from 'vuex';

const axios = require('axios');

const client = axios.create({
  baseURL: 'http://localhost:3000/',
});

Vue.use(Vuex);

/*
const FakeVideoListing = [
  {
    id: 1,
    title: 'dinotown',
    description: 'dinosaurs',
    thumbnail: 'http://localhost:3000/static/data-xor/dinotown.mp4.jpg',
    source: 'http://localhost:3000/static/data-xor/dinotown.mp4',
    time_created: '2018-12-25T00:00:00Z',
    time_updated: '2018-12-25T00:00:00Z',
    tags: [
      'nostalgia',
    ],
  },
  {
    id: 2,
    title: 'meme',
    description: 'a spicy maymay',
    thumbnail: 'http://localhost:3000/static/data-xor/meme.webm.jpg',
    source: 'http://localhost:3000/static/data-xor/meme.webm',
    time_created: '2018-12-25T00:00:00Z',
    time_updated: '2018-12-25T00:00:00Z',
    tags: [
      'meme',
      'spicy',
    ],
  },
  {
    id: 3,
    title: 'Last Chance To Evacuate Earth',
    description: 'Recorded September 29, 1996.',
    thumbnail: 'http://localhost:3000/static/data-xor/Last Chance To Evacuate Earth Before It\'s Recycled.mp4.jpg',
    source: 'http://localhost:3000/static/data-xor/Last Chance To Evacuate Earth Before It\'s Recycled.mp4',
    time_created: '2018-12-25T00:00:00Z',
    time_updated: '2018-12-25T00:00:00Z',
    tags: [
      'heavens gate',
    ],
  },
  {
    id: 4,
    title: 'Seagulls',
    description: 'Some seagull stuff',
    thumbnail: 'http://localhost:3000/static/data-xor/seagulls.mp4.jpg',
    source: 'http://localhost:3000/static/data-xor/seagulls.mp4',
    time_created: '2018-12-25T00:00:00Z',
    time_updated: '2018-12-25T00:00:00Z',
    tags: [
      'meme',
    ],
  },
];

const fakePromiseDelay = (delay = 0) => new Promise((resolve) => {
  setTimeout(() => resolve(), delay);
});
*/

export default new Vuex.Store({
  state: {
    videos: [],
  },
  mutations: {
    setVideos(state, videos) {
      state.videos = videos; // eslint-disable-line
    },
  },
  actions: {
    videos({ commit }) {
      return client.get('/api/video')
        .then(resp => resp.data)
        .then((videos) => {
          commit('setVideos', videos);
          return videos;
        });
    },
    video({ dispatch }, id) {
      // todo: replace with actual /api/video/{id}
      return dispatch('videos').then(videos => videos.find(v => v.id === id));
    },
    upload(context, formData) {
      return client.post('/api/upload', formData);
    },
  },
});

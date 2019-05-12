<template>
  <select v-model="sortOption">
    <option
      v-for="option in sortOptions"
      :key="option.key"
      :value="option"
    >Sort: {{ option.name }}</option>
  </select>
</template>

<script>
import sortOptions from '@/sortOptions';

export default {
  computed: {
    sortOptions() {
      return sortOptions;
    },

    sortOption: {
      get() {
        return this.sortOptions.find(option => option.key === this.value) || this.sortOptions[0];
      },

      set(value) {
        this.$emit('input', value.key);
      },
    },
  },

  props: {
    value: {
      type: Object,
      required: false,
      default: () => {
        return sortOptions[0].key;
      },
    },  
  },
};
</script>

<style lang="scss" scoped>
select {
  text-align: right;
  appearance: none;
  color: white;
  border: 0;
  line-height: 1;
  outline: 0;
  background-color: transparent;
  font-size: 1em;
  color: rgba(255,255,255,0.6);

  option {
    color: black;
    background-color: white;
  }
}
</style>

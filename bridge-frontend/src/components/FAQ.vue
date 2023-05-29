<template>
  <div class="faq-container">
    <h2>Frequently Asked Questions</h2>
    <ul>
      <li v-for="entry in faqEntries" :key="entry.question">
        <div class="question-block" @click="toggleAnswer(entry)">
          <h3>{{ entry.question }}</h3>
          <p v-show="entry.showAnswer">{{ entry.answer }}</p>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: 'FAQ',
  data() {
    return {
      faqEntries: [
        {
          question: 'What is Specular?',
          answer: 'Specular Network is a Layer 2 scaling solution for Ethereum that uses optimistic rollup technology to scale the Ethereum network. It is still under development, but it has the potential to provide a number of benefits for Ethereum users, including scalability, security, and decentralization.',
          showAnswer: false,
        },
        {
          question: 'What is Specular Bridge?',
          answer: 'The Specular Network Bridge is a two-way bridge that allows users to transfer assets between Ethereum and Specular Network. It is still in beta, but it is expected to be released to the public in the near future. The bridge offers low fees, high throughput, and security.',
          showAnswer: false,
        },
        {
          question: 'What chains are supported By Specular Bridge?',
          answer: 'As of today, Specular Network only supports Chiado. However, the team has plans to support other chains in the future.',
          showAnswer: false,
        },
      ],
    };
  },
  methods: {
    toggleAnswer(entry) {
      entry.showAnswer = !entry.showAnswer;
    },
  },
};
</script>

<style scoped>
.faq-container {
  padding-top: 50px;
  padding-bottom: 10px;
}

h2 {
  font-size: 24px;
  margin-bottom: 20px;
}

.question-block {
  background-color: #f2f2f2;
  padding: 10px;
  margin-bottom: 10px;
  cursor: pointer;
  border-radius: 8px;
}

h3 {
  font-size: 18px;
  margin-bottom: 5px;
}

p {
  font-size: 14px;
  margin-bottom: 10px;
}
</style>

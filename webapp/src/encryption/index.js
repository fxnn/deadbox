export default {
  state: {
    keyAvailable: false
  },
  actions: {
    setKeyAvailable: () => (state) => ({ keyAvailable: true }),
    setKeyAvailableDelayed: timeout => (state, actions) =>
      new Promise(resolve => setTimeout(resolve, timeout)).then(() => actions.setKeyAvailable())
  }
};
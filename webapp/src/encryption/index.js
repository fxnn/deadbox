export default {
  state: {
    keyAvailable: false,
    keyConfigurationAvailable: false
  },
  actions: {
    setKeyConfigured: () => (state) => ({ keyConfigurationAvailable: true }),
    setKeyAvailable: () => (state) => ({ keyAvailable: true }),
    setKeyAvailableDelayed: timeout => (state, actions) =>
      new Promise(resolve => setTimeout(resolve, timeout)).then(() => actions.setKeyAvailable())
  }
};
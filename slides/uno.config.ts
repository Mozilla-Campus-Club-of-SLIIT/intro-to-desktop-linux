import { defineConfig } from "unocss";

export default defineConfig({
  theme: {
    colors: {
      black: "#000000",
      white: "#FFFFFF",
      lemonYellow: "#fff44f",
      warmRed: "#ff4f5e",
      darkGreen: "#005e5e",
    },
  },
  shortcuts: {
    // This overrides the default Slidev background and text color globally
    "bg-main": "bg-black text-white",
  },
});

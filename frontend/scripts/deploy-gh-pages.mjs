import { publish } from "gh-pages";

const branch = process.env.GH_PAGES_BRANCH || "gh-pages";

publish(
  "dist",
  {
    branch,
    dotfiles: true,
    message: "deploy frontend to GitHub Pages",
  },
  (error) => {
    if (error) {
      console.error(error);
      process.exit(1);
    }

    console.log(`published dist to ${branch}`);
  },
);

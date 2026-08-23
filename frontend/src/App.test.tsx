import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { App } from "./App";

describe("App", () => {
  it("renders the login screen by default when unauthenticated", async () => {
    render(<App />);

    expect(await screen.findByText("Bem-vindo de volta")).toBeInTheDocument();
    expect(
      await screen.findByRole("button", { name: /Entrar no CVMC/i }),
    ).toBeInTheDocument();
    expect(
      await screen.findByPlaceholderText("seu.email@exemplo.com"),
    ).toBeInTheDocument();
  });
});

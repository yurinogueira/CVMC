import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { App } from "./App";

describe("App", () => {
  it("renders the login screen by default when unauthenticated", () => {
    render(<App />);

    expect(screen.getByText("Bem-vindo de volta")).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /Entrar no CVMC/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByPlaceholderText("seu.email@exemplo.com"),
    ).toBeInTheDocument();
  });
});

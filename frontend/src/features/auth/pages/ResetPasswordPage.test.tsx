import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { ResetPasswordPage } from "./ResetPasswordPage";

describe("ResetPasswordPage", () => {
  it("renders reset password form when token is provided in query params", () => {
    render(
      <MemoryRouter initialEntries={["/reset-password?token=valid-token-123"]}>
        <ResetPasswordPage />
      </MemoryRouter>,
    );

    expect(
      screen.getByRole("heading", { name: /Nova Senha/i }),
    ).toBeInTheDocument();
    expect(screen.getByLabelText(/^Nova Senha/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/^Confirmar Nova Senha/i)).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /Salvar Nova Senha/i }),
    ).toBeInTheDocument();
  });

  it("shows warning alert when token is missing", () => {
    render(
      <MemoryRouter initialEntries={["/reset-password"]}>
        <ResetPasswordPage />
      </MemoryRouter>,
    );

    expect(
      screen.getByText(/Token de recuperação não identificado/i),
    ).toBeInTheDocument();
  });
});

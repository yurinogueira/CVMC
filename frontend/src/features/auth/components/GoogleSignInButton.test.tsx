import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { GoogleSignInButton } from "./GoogleSignInButton";
import { Capacitor } from "@capacitor/core";
import { SocialLogin } from "@capgo/capacitor-social-login";

vi.mock("@capacitor/core", () => ({
  Capacitor: {
    isNativePlatform: vi.fn(),
  },
}));

vi.mock("@capgo/capacitor-social-login", () => ({
  SocialLogin: {
    initialize: vi.fn().mockResolvedValue(undefined),
    login: vi.fn(),
  },
}));

describe("GoogleSignInButton", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders web fallback button when not in native environment", () => {
    vi.mocked(Capacitor.isNativePlatform).mockReturnValue(false);

    render(
      <GoogleSignInButton
        onSuccess={vi.fn()}
        text="signin_with"
        label="Entrar com Google"
      />,
    );

    expect(screen.getByTestId("google-signin-btn")).toBeInTheDocument();
    expect(screen.getByText("Entrar com Google")).toBeInTheDocument();
  });

  it("renders custom label when provided", () => {
    vi.mocked(Capacitor.isNativePlatform).mockReturnValue(false);

    render(
      <GoogleSignInButton onSuccess={vi.fn()} label="Vincular Conta Google" />,
    );

    expect(screen.getByText("Vincular Conta Google")).toBeInTheDocument();
  });

  it("triggers native SocialLogin on click when isNativePlatform is true and succeeds", async () => {
    vi.mocked(Capacitor.isNativePlatform).mockReturnValue(true);
    const onSuccess = vi.fn();
    const onError = vi.fn();

    vi.mocked(SocialLogin.login).mockResolvedValueOnce({
      provider: "google",
      result: {
        idToken: "mocked-native-google-id-token",
        accessToken: null,
        profile: {
          email: "user@example.com",
          familyName: "User",
          givenName: "Test",
          id: "123",
          name: "Test User",
          imageUrl: null,
        },
        responseType: "online",
      },
    });

    render(
      <GoogleSignInButton
        onSuccess={onSuccess}
        onError={onError}
        label="Entrar com Google"
      />,
    );

    const button = screen.getByTestId("google-signin-btn");
    fireEvent.click(button);

    await waitFor(() => {
      expect(SocialLogin.initialize).toHaveBeenCalled();
      expect(SocialLogin.login).toHaveBeenCalledWith({
        provider: "google",
        options: {
          filterByAuthorizedAccounts: false,
        },
      });
      expect(onSuccess).toHaveBeenCalledWith("mocked-native-google-id-token");
      expect(onError).not.toHaveBeenCalled();
    });
  });

  it("handles user cancellation gracefully without calling onError", async () => {
    vi.mocked(Capacitor.isNativePlatform).mockReturnValue(true);
    const onSuccess = vi.fn();
    const onError = vi.fn();

    vi.mocked(SocialLogin.login).mockRejectedValueOnce({
      code: "USER_CANCELLED",
      message: "User cancelled login dialog",
    });

    render(
      <GoogleSignInButton
        onSuccess={onSuccess}
        onError={onError}
        label="Entrar com Google"
      />,
    );

    const button = screen.getByTestId("google-signin-btn");
    fireEvent.click(button);

    await waitFor(() => {
      expect(SocialLogin.login).toHaveBeenCalled();
      expect(onSuccess).not.toHaveBeenCalled();
      expect(onError).not.toHaveBeenCalled();
    });
  });

  it("calls onError when native login fails with an error", async () => {
    vi.mocked(Capacitor.isNativePlatform).mockReturnValue(true);
    const onSuccess = vi.fn();
    const onError = vi.fn();

    vi.mocked(SocialLogin.login).mockRejectedValueOnce(
      new Error("Play Services unavailable"),
    );

    render(
      <GoogleSignInButton
        onSuccess={onSuccess}
        onError={onError}
        label="Entrar com Google"
      />,
    );

    const button = screen.getByTestId("google-signin-btn");
    fireEvent.click(button);

    await waitFor(() => {
      expect(onError).toHaveBeenCalledWith(expect.any(Error));
      expect(onSuccess).not.toHaveBeenCalled();
    });
  });
});

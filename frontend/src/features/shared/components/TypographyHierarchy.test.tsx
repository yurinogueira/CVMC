import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { PageHeader } from "./PageHeader";
import { getTopbarTitle, staticRouteTitles } from "../../../layouts/Topbar";

describe("Typography Hierarchy & Visual Consistency Architecture Tests", () => {
  it("guarantees PageHeader enforces single semantic h1 tag with variant h5 typography", () => {
    render(
      <PageHeader
        title="Painel de Controle"
        subtitle="Métricas operacionais"
      />,
    );

    const heading = screen.getByRole("heading", {
      level: 1,
      name: "Painel de Controle",
    });
    expect(heading).toBeInTheDocument();
    expect(heading.tagName).toBe("H1");

    // Must be exactly variant="h5" in MUI classes (MuiTypography-h5)
    expect(heading.className).toContain("MuiTypography-h5");

    const subtitle = screen.getByText("Métricas operacionais");
    expect(subtitle.className).toContain("MuiTypography-body2");
  });

  it("guarantees all mandatory application routes have explicit non-fallback Topbar titles", () => {
    const protectedRoutes = [
      "/dashboard",
      "/vehicles",
      "/vehicles/car-test-1",
      "/vehicles/car-test-1/maintenance/new",
      "/vehicles/car-test-1/maintenance/m-99/edit",
      "/maintenance",
      "/reports",
      "/profile",
      "/privacy",
    ];

    for (const route of protectedRoutes) {
      const title = getTopbarTitle(route);
      expect(title).not.toBe("Painel CVMC");
      expect(title.length).toBeGreaterThan(0);
    }
  });

  it("guarantees static route titles map exactly to canonical page headers", () => {
    expect(staticRouteTitles["/dashboard"]).toBe("Dashboard");
    expect(staticRouteTitles["/vehicles"]).toBe("Meus Veículos");
    expect(staticRouteTitles["/maintenance"]).toBe("Manutenções & Revisões");
    expect(staticRouteTitles["/reports"]).toBe("Relatórios & Exportações");
    expect(staticRouteTitles["/profile"]).toBe("Meu Perfil");
    expect(staticRouteTitles["/privacy"]).toBe("Política de Privacidade");
  });

  it("guarantees trailing slashes on all routes resolve without falling back", () => {
    const routesWithTrailingSlash = [
      "/dashboard/",
      "/vehicles/",
      "/maintenance/",
      "/reports/",
      "/profile/",
    ];

    for (const route of routesWithTrailingSlash) {
      const title = getTopbarTitle(route);
      expect(title).not.toBe("Painel CVMC");
    }
  });
});

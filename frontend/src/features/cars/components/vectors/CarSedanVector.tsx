import React from "react";
import { Box, BoxProps } from "@mui/material";

export function CarSedanVector(props: BoxProps) {
  return (
    <Box
      component="svg"
      viewBox="0 0 800 450"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      sx={{
        width: "100%",
        height: "100%",
        display: "block",
        ...props.sx,
      }}
      {...props}
    >
      <defs>
        {/* Sky/Atmosphere Backdrop */}
        <radialGradient id="carBgGlow" cx="50%" cy="40%" r="60%">
          <stop offset="0%" stopColor="#38BDF8" stopOpacity="0.15" />
          <stop offset="60%" stopColor="#0284C7" stopOpacity="0.05" />
          <stop offset="100%" stopColor="#0F172A" stopOpacity="0" />
        </radialGradient>

        {/* Car Metallic Paint Gradient */}
        <linearGradient
          id="carBodyGrad"
          x1="150"
          y1="120"
          x2="680"
          y2="350"
          gradientUnits="userSpaceOnUse"
        >
          <stop offset="0%" stopColor="#38BDF8" />
          <stop offset="25%" stopColor="#0284C7" />
          <stop offset="60%" stopColor="#0369A1" />
          <stop offset="100%" stopColor="#0C4A6E" />
        </linearGradient>

        {/* Highlight Reflection */}
        <linearGradient
          id="carHighlight"
          x1="200"
          y1="120"
          x2="550"
          y2="220"
          gradientUnits="userSpaceOnUse"
        >
          <stop offset="0%" stopColor="#FFFFFF" stopOpacity="0.6" />
          <stop offset="30%" stopColor="#BAE6FD" stopOpacity="0.2" />
          <stop offset="100%" stopColor="#FFFFFF" stopOpacity="0" />
        </linearGradient>

        {/* Glass Tint */}
        <linearGradient
          id="glassGrad"
          x1="280"
          y1="140"
          x2="540"
          y2="240"
          gradientUnits="userSpaceOnUse"
        >
          <stop offset="0%" stopColor="#0F172A" stopOpacity="0.95" />
          <stop offset="40%" stopColor="#1E293B" stopOpacity="0.85" />
          <stop offset="100%" stopColor="#38BDF8" stopOpacity="0.3" />
        </linearGradient>

        {/* Wheel Rims Gradient */}
        <radialGradient id="rimGrad" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#E2E8F0" />
          <stop offset="60%" stopColor="#94A3B8" />
          <stop offset="85%" stopColor="#334155" />
          <stop offset="100%" stopColor="#0F172A" />
        </radialGradient>

        {/* Headlight LED Beam */}
        <linearGradient
          id="ledGlow"
          x1="680"
          y1="260"
          x2="780"
          y2="280"
          gradientUnits="userSpaceOnUse"
        >
          <stop offset="0%" stopColor="#38BDF8" stopOpacity="0.9" />
          <stop offset="60%" stopColor="#38BDF8" stopOpacity="0.3" />
          <stop offset="100%" stopColor="#38BDF8" stopOpacity="0" />
        </linearGradient>

        {/* Ground Shadow */}
        <radialGradient id="shadowGrad" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#0F172A" stopOpacity="0.6" />
          <stop offset="60%" stopColor="#0F172A" stopOpacity="0.3" />
          <stop offset="100%" stopColor="#0F172A" stopOpacity="0" />
        </radialGradient>
      </defs>

      {/* Subtle Background Glow */}
      <rect width="800" height="450" fill="url(#carBgGlow)" />

      {/* Ground Contact Shadow */}
      <ellipse cx="440" cy="365" rx="330" ry="32" fill="url(#shadowGrad)" />
      <ellipse cx="270" cy="362" rx="75" ry="18" fill="#000000" opacity="0.6" />
      <ellipse cx="585" cy="362" rx="85" ry="18" fill="#000000" opacity="0.6" />

      {/* Headlight Projection Glow */}
      <polygon points="660,265 780,250 790,320 670,300" fill="url(#ledGlow)" />

      {/* CAR BODY - 3/4 Perspective Sedan */}
      {/* Rear & Lower Undercarriage */}
      <path
        d="M130 310 C140 280, 180 270, 220 270 L240 270 C250 250, 280 235, 320 235 L530 235 C565 235, 595 250, 610 275 L675 290 C700 295, 715 315, 710 330 C705 345, 680 355, 640 355 L160 355 C135 355, 120 335, 130 310 Z"
        fill="#075985"
      />

      {/* Main Aerodynamic Shell */}
      <path
        d="M140 310 
           C150 255, 200 235, 260 215 
           L340 145 
           C365 125, 430 120, 510 135 
           L595 195 
           C635 225, 680 250, 710 285 
           C720 298, 715 318, 695 328 
           L675 338 
           C660 345, 620 348, 560 348 
           L240 348 
           C180 348, 140 335, 140 310 Z"
        fill="url(#carBodyGrad)"
      />

      {/* Cabin Roofline & Greenhouse */}
      <path
        d="M260 215 
           L345 145 
           C370 125, 435 120, 510 135 
           L595 195 
           C600 200, 580 215, 530 220 
           L330 220 
           C290 220, 265 218, 260 215 Z"
        fill="url(#glassGrad)"
        stroke="#0284C7"
        strokeWidth="2"
      />

      {/* Window Pillar (B-Pillar) */}
      <polygon points="420,132 435,133 430,218 415,218" fill="#0F172A" />

      {/* Windshield Reflection Streak */}
      <path
        d="M355 148 L500 138 L430 212 L300 212 Z"
        fill="url(#carHighlight)"
        opacity="0.7"
      />

      {/* Side Windows (Front & Rear) */}
      <path
        d="M335 210 L355 155 C375 140, 410 136, 415 136 L410 210 Z"
        fill="#0F172A"
        opacity="0.7"
      />
      <path
        d="M440 210 L445 137 C480 142, 530 165, 565 198 L510 210 Z"
        fill="#0F172A"
        opacity="0.65"
      />

      {/* Dynamic Shoulder Swage Line */}
      <path
        d="M150 295 C220 270, 360 250, 540 255 C620 258, 675 280, 710 290"
        stroke="url(#carHighlight)"
        strokeWidth="3.5"
        strokeLinecap="round"
      />

      {/* Hood Character Lines */}
      <path
        d="M510 140 L640 235 L700 285"
        stroke="#38BDF8"
        strokeWidth="2"
        opacity="0.8"
      />
      <path
        d="M595 195 L690 265"
        stroke="#7DD3FC"
        strokeWidth="1.5"
        opacity="0.6"
      />

      {/* Side Mirror */}
      <path
        d="M510 195 C525 190, 535 200, 530 210 C520 215, 505 210, 505 200 Z"
        fill="#0369A1"
        stroke="#38BDF8"
        strokeWidth="1.5"
      />

      {/* Door Seam Lines */}
      <path
        d="M425 218 L420 340"
        stroke="#0369A1"
        strokeWidth="2"
        opacity="0.8"
      />
      <path
        d="M315 220 L305 340"
        stroke="#0369A1"
        strokeWidth="2"
        opacity="0.8"
      />
      <path
        d="M525 220 L515 340"
        stroke="#0369A1"
        strokeWidth="2"
        opacity="0.8"
      />

      {/* Chrome Door Handles */}
      <rect x="355" y="245" width="22" height="5" rx="2.5" fill="#E2E8F0" />
      <rect x="460" y="245" width="22" height="5" rx="2.5" fill="#E2E8F0" />

      {/* Front Headlight (LED Projector 3/4) */}
      <path
        d="M660 270 C685 275, 705 288, 712 300 C705 308, 680 305, 655 295 Z"
        fill="#E0F2FE"
        stroke="#38BDF8"
        strokeWidth="2"
      />
      <circle cx="680" cy="288" r="6" fill="#38BDF8" />
      <path
        d="M665 278 L700 292"
        stroke="#FFFFFF"
        strokeWidth="2.5"
        strokeLinecap="round"
      />

      {/* Front Grille & Lower Bumper Air Dam */}
      <path
        d="M685 305 L705 320 C695 335, 670 342, 635 345 L625 330 Z"
        fill="#0F172A"
      />
      <line
        x1="640"
        y1="335"
        x2="685"
        y2="320"
        stroke="#334155"
        strokeWidth="2"
      />
      <line
        x1="645"
        y1="340"
        x2="690"
        y2="325"
        stroke="#334155"
        strokeWidth="2"
      />

      {/* Rear Taillight Strip */}
      <path
        d="M140 298 C145 290, 160 288, 175 290 L168 308 Z"
        fill="#EF4444"
        stroke="#F87171"
        strokeWidth="1"
      />

      {/* REAR WHEEL ASSEMBLY */}
      <g>
        {/* Wheel Arch Cavity */}
        <path
          d="M205 348 C205 305, 240 275, 285 275 C330 275, 360 305, 360 348 Z"
          fill="#0F172A"
        />
        {/* Tire */}
        <circle
          cx="280"
          cy="345"
          r="54"
          fill="#1E293B"
          stroke="#0F172A"
          strokeWidth="4"
        />
        <circle cx="280" cy="345" r="40" fill="url(#rimGrad)" />
        {/* Alloy Spokes (Sedan Sport Design) */}
        <circle cx="280" cy="345" r="14" fill="#0F172A" />
        <line
          x1="280"
          y1="310"
          x2="280"
          y2="380"
          stroke="#F8FAFC"
          strokeWidth="4.5"
          strokeLinecap="round"
        />
        <line
          x1="245"
          y1="345"
          x2="315"
          y2="345"
          stroke="#F8FAFC"
          strokeWidth="4.5"
          strokeLinecap="round"
        />
        <line
          x1="255"
          y1="320"
          x2="305"
          y2="370"
          stroke="#CBD5E1"
          strokeWidth="3.5"
          strokeLinecap="round"
        />
        <line
          x1="255"
          y1="370"
          x2="305"
          y2="320"
          stroke="#CBD5E1"
          strokeWidth="3.5"
          strokeLinecap="round"
        />
        <circle cx="280" cy="345" r="6" fill="#38BDF8" />
      </g>

      {/* FRONT WHEEL ASSEMBLY (Larger due to perspective) */}
      <g>
        {/* Wheel Arch Cavity */}
        <path
          d="M515 348 C515 295, 555 260, 610 260 C660 260, 695 295, 695 348 Z"
          fill="#0F172A"
        />
        {/* Tire */}
        <circle
          cx="600"
          cy="345"
          r="60"
          fill="#1E293B"
          stroke="#0F172A"
          strokeWidth="5"
        />
        <circle cx="600" cy="345" r="45" fill="url(#rimGrad)" />
        {/* Alloy Spokes */}
        <circle cx="600" cy="345" r="16" fill="#0F172A" />
        <line
          x1="600"
          y1="305"
          x2="600"
          y2="385"
          stroke="#F8FAFC"
          strokeWidth="5"
          strokeLinecap="round"
        />
        <line
          x1="560"
          y1="345"
          x2="640"
          y2="345"
          stroke="#F8FAFC"
          strokeWidth="5"
          strokeLinecap="round"
        />
        <line
          x1="570"
          y1="315"
          x2="630"
          y2="375"
          stroke="#CBD5E1"
          strokeWidth="4"
          strokeLinecap="round"
        />
        <line
          x1="570"
          y1="375"
          x2="630"
          y2="315"
          stroke="#CBD5E1"
          strokeWidth="4"
          strokeLinecap="round"
        />
        {/* Brake Caliper (Performance Blue) */}
        <path
          d="M625 320 C635 330, 635 350, 625 360"
          stroke="#0284C7"
          strokeWidth="6"
          strokeLinecap="round"
          fill="none"
        />
        <circle cx="600" cy="345" r="7" fill="#38BDF8" />
      </g>
    </Box>
  );
}

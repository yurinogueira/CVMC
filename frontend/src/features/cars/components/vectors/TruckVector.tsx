import React from "react";
import { Box, BoxProps } from "@mui/material";

export function TruckVector(props: BoxProps) {
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
        {/* Background Glow */}
        <radialGradient id="truckBgGlow" cx="50%" cy="40%" r="60%">
          <stop offset="0%" stopColor="#38BDF8" stopOpacity="0.15" />
          <stop offset="60%" stopColor="#0284C7" stopOpacity="0.05" />
          <stop offset="100%" stopColor="#0F172A" stopOpacity="0" />
        </radialGradient>

        {/* Truck Metallic Gradient */}
        <linearGradient
          id="truckBodyGrad"
          x1="120"
          y1="100"
          x2="680"
          y2="360"
          gradientUnits="userSpaceOnUse"
        >
          <stop offset="0%" stopColor="#38BDF8" />
          <stop offset="30%" stopColor="#0284C7" />
          <stop offset="70%" stopColor="#075985" />
          <stop offset="100%" stopColor="#0F172A" />
        </linearGradient>

        {/* Highlight Reflection */}
        <linearGradient
          id="truckHighlight"
          x1="160"
          y1="120"
          x2="580"
          y2="240"
          gradientUnits="userSpaceOnUse"
        >
          <stop offset="0%" stopColor="#FFFFFF" stopOpacity="0.6" />
          <stop offset="50%" stopColor="#BAE6FD" stopOpacity="0.25" />
          <stop offset="100%" stopColor="#FFFFFF" stopOpacity="0" />
        </linearGradient>

        {/* Tinted Windows */}
        <linearGradient
          id="truckGlass"
          x1="300"
          y1="120"
          x2="560"
          y2="220"
          gradientUnits="userSpaceOnUse"
        >
          <stop offset="0%" stopColor="#0F172A" stopOpacity="0.95" />
          <stop offset="60%" stopColor="#1E293B" stopOpacity="0.85" />
          <stop offset="100%" stopColor="#38BDF8" stopOpacity="0.25" />
        </linearGradient>

        {/* Offroad Heavy Rim Gradient */}
        <radialGradient id="truckRim" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#F1F5F9" />
          <stop offset="50%" stopColor="#64748B" />
          <stop offset="85%" stopColor="#1E293B" />
          <stop offset="100%" stopColor="#0F172A" />
        </radialGradient>

        {/* Headlight LED Beam */}
        <linearGradient
          id="truckLedGlow"
          x1="680"
          y1="230"
          x2="790"
          y2="260"
          gradientUnits="userSpaceOnUse"
        >
          <stop offset="0%" stopColor="#38BDF8" stopOpacity="0.9" />
          <stop offset="100%" stopColor="#38BDF8" stopOpacity="0" />
        </linearGradient>

        {/* Ground Shadow */}
        <radialGradient id="truckShadow" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#0F172A" stopOpacity="0.7" />
          <stop offset="70%" stopColor="#0F172A" stopOpacity="0.25" />
          <stop offset="100%" stopColor="#0F172A" stopOpacity="0" />
        </radialGradient>
      </defs>

      {/* Background Glow */}
      <rect width="800" height="450" fill="url(#truckBgGlow)" />

      {/* Ground Contact Shadow */}
      <ellipse cx="440" cy="375" rx="340" ry="32" fill="url(#truckShadow)" />
      <ellipse
        cx="260"
        cy="370"
        rx="85"
        ry="20"
        fill="#000000"
        opacity="0.65"
      />
      <ellipse
        cx="600"
        cy="370"
        rx="95"
        ry="20"
        fill="#000000"
        opacity="0.65"
      />

      {/* Headlight Beam */}
      <polygon
        points="680,225 790,205 800,310 690,290"
        fill="url(#truckLedGlow)"
      />

      {/* TRUCK CABIN & BED SILHOUETTE (3/4 Pick-up / Truck) */}
      {/* Lower chassis / side step rail */}
      <rect
        x="290"
        y="340"
        width="260"
        height="12"
        rx="4"
        fill="#0F172A"
        stroke="#334155"
        strokeWidth="2"
      />

      {/* Main Truck Body */}
      <path
        d="M100 240 
           L260 235 
           L300 135 
           C315 110, 360 100, 480 105 
           L560 180 
           L680 200 
           C705 205, 725 225, 725 255 
           L715 325 
           C705 345, 680 350, 620 350 
           L130 350 
           C105 350, 95 330, 95 290 Z"
        fill="url(#truckBodyGrad)"
      />

      {/* Pick-up Bed Cargo Area Line & Rollbar */}
      <polygon
        points="100,240 260,235 260,250 105,255"
        fill="#0F172A"
        opacity="0.7"
      />
      {/* Sport Rollbar / Bed Rails */}
      <path
        d="M250 235 L285 140 L310 140 L280 235"
        fill="#334155"
        stroke="#1E293B"
        strokeWidth="2"
      />
      <path
        d="M285 140 L370 140"
        stroke="#334155"
        strokeWidth="6"
        strokeLinecap="round"
      />

      {/* Cabin Windows (Heavy Duty Dual Cab Glass) */}
      <path
        d="M300 140 
           C315 115, 360 108, 475 112 
           L550 180 
           L520 215 
           L310 215 Z"
        fill="url(#truckGlass)"
        stroke="#0284C7"
        strokeWidth="2"
      />
      {/* Center Pillar */}
      <polygon points="410,110 425,110 420,215 405,215" fill="#0F172A" />

      {/* Windshield Reflection */}
      <path
        d="M330 145 L465 120 L530 180 L420 210 Z"
        fill="url(#truckHighlight)"
        opacity="0.6"
      />

      {/* Side Windows */}
      <path
        d="M320 208 L330 148 C345 128, 395 120, 400 120 L395 208 Z"
        fill="#0F172A"
        opacity="0.7"
      />
      <path
        d="M430 208 L435 122 C465 128, 515 155, 535 185 L495 208 Z"
        fill="#0F172A"
        opacity="0.65"
      />

      {/* Heavy-Duty Side Mirror */}
      <path
        d="M510 185 C528 180, 540 190, 535 205 C525 212, 505 208, 505 195 Z"
        fill="#0F172A"
        stroke="#38BDF8"
        strokeWidth="2"
      />

      {/* Muscular Shoulder & Fender Flares */}
      <path
        d="M100 280 C200 260, 360 245, 540 250 C630 255, 680 270, 720 275"
        stroke="url(#truckHighlight)"
        strokeWidth="4"
        strokeLinecap="round"
      />

      {/* Hood Power Bulge Lines */}
      <path d="M475 112 L600 205 L690 230" stroke="#38BDF8" strokeWidth="2.5" />
      <path
        d="M555 180 L675 220"
        stroke="#7DD3FC"
        strokeWidth="2"
        opacity="0.7"
      />

      {/* Door Seam Lines */}
      <path d="M415 215 L405 340" stroke="#075985" strokeWidth="2.5" />
      <path d="M295 235 L285 340" stroke="#075985" strokeWidth="2.5" />
      <path d="M515 215 L505 340" stroke="#075985" strokeWidth="2.5" />

      {/* Rugged Door Handles */}
      <rect
        x="345"
        y="240"
        width="24"
        height="7"
        rx="3"
        fill="#0F172A"
        stroke="#64748B"
        strokeWidth="1.5"
      />
      <rect
        x="445"
        y="240"
        width="24"
        height="7"
        rx="3"
        fill="#0F172A"
        stroke="#64748B"
        strokeWidth="1.5"
      />

      {/* IMPOSING FRONT GRILLE & BUMPER (Chrome/Blackout RAM/Hilux Style) */}
      <path
        d="M660 215 
           L720 230 
           L715 315 
           L650 335 
           L640 250 Z"
        fill="#0F172A"
        stroke="#334155"
        strokeWidth="2"
      />
      {/* Grille Bars */}
      <line
        x1="660"
        y1="245"
        x2="710"
        y2="255"
        stroke="#38BDF8"
        strokeWidth="3"
      />
      <line
        x1="655"
        y1="265"
        x2="705"
        y2="275"
        stroke="#38BDF8"
        strokeWidth="3"
      />
      <line
        x1="650"
        y1="285"
        x2="700"
        y2="295"
        stroke="#38BDF8"
        strokeWidth="3"
      />

      {/* Front Skid Plate / Bull Bar */}
      <path
        d="M640 325 L710 310 L705 345 L635 348 Z"
        fill="#334155"
        stroke="#64748B"
        strokeWidth="2"
      />

      {/* Front Headlight (Aggressive C-Clamp LED) */}
      <path
        d="M675 220 C695 225, 715 235, 720 255 C710 265, 685 260, 670 248 Z"
        fill="#E0F2FE"
        stroke="#38BDF8"
        strokeWidth="2.5"
      />
      <circle cx="695" cy="242" r="7" fill="#38BDF8" />
      <path
        d="M680 228 L715 242"
        stroke="#FFFFFF"
        strokeWidth="3"
        strokeLinecap="round"
      />

      {/* Rear Taillight Cluster */}
      <path
        d="M98 245 L115 245 L110 295 L95 295 Z"
        fill="#EF4444"
        stroke="#F87171"
        strokeWidth="1.5"
      />

      {/* REAR WHEEL ASSEMBLY (Heavy Offroad Tire + Arch Flare) */}
      <g>
        {/* Bulging Wheel Flare */}
        <path
          d="M190 350 C190 295, 230 260, 285 260 C335 260, 365 295, 365 350 Z"
          fill="#0F172A"
        />
        {/* Offroad Knobby Tire */}
        <circle
          cx="275"
          cy="345"
          r="66"
          fill="#1E293B"
          stroke="#0F172A"
          strokeWidth="6"
        />
        <circle cx="275" cy="345" r="46" fill="url(#truckRim)" />
        {/* 6-Lug Truck Rim Design */}
        <circle cx="275" cy="345" r="18" fill="#0F172A" />
        <line
          x1="275"
          y1="305"
          x2="275"
          y2="385"
          stroke="#F8FAFC"
          strokeWidth="6"
          strokeLinecap="round"
        />
        <line
          x1="235"
          y1="345"
          x2="315"
          y2="345"
          stroke="#F8FAFC"
          strokeWidth="6"
          strokeLinecap="round"
        />
        <line
          x1="245"
          y1="315"
          x2="305"
          y2="375"
          stroke="#94A3B8"
          strokeWidth="5"
          strokeLinecap="round"
        />
        <line
          x1="245"
          y1="375"
          x2="305"
          y2="315"
          stroke="#94A3B8"
          strokeWidth="5"
          strokeLinecap="round"
        />
        {/* Center Cap */}
        <circle cx="275" cy="345" r="8" fill="#38BDF8" />
      </g>

      {/* FRONT WHEEL ASSEMBLY (Larger in 3/4 perspective) */}
      <g>
        {/* Bulging Wheel Flare */}
        <path
          d="M510 350 C510 285, 555 245, 620 245 C675 245, 710 285, 710 350 Z"
          fill="#0F172A"
        />
        {/* Offroad Knobby Tire */}
        <circle
          cx="610"
          cy="345"
          r="74"
          fill="#1E293B"
          stroke="#0F172A"
          strokeWidth="6"
        />
        <circle cx="610" cy="345" r="52" fill="url(#truckRim)" />
        {/* 6-Lug Truck Rim */}
        <circle cx="610" cy="345" r="20" fill="#0F172A" />
        <line
          x1="610"
          y1="300"
          x2="610"
          y2="390"
          stroke="#F8FAFC"
          strokeWidth="7"
          strokeLinecap="round"
        />
        <line
          x1="565"
          y1="345"
          x2="655"
          y2="345"
          stroke="#F8FAFC"
          strokeWidth="7"
          strokeLinecap="round"
        />
        <line
          x1="575"
          y1="310"
          x2="645"
          y2="380"
          stroke="#94A3B8"
          strokeWidth="6"
          strokeLinecap="round"
        />
        <line
          x1="575"
          y1="380"
          x2="645"
          y2="310"
          stroke="#94A3B8"
          strokeWidth="6"
          strokeLinecap="round"
        />
        {/* Brake Caliper */}
        <path
          d="M635 320 C648 332, 648 358, 635 370"
          stroke="#0284C7"
          strokeWidth="7"
          strokeLinecap="round"
          fill="none"
        />
        <circle cx="610" cy="345" r="9" fill="#38BDF8" />
      </g>
    </Box>
  );
}

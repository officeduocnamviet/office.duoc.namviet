"use client";

import React, { useState } from "react";
import Link from "next/link";
import { Grid, ShoppingBag, ShoppingCart, User } from "lucide-react";

interface MobileBottomNavProps {
  activeTab: string;
  setActiveTab: (tab: string) => void;
  cartCount: number;
  isFabOpen: boolean;
  toggleFabMenu: () => void;
}

export default function MobileBottomNav({
  activeTab,
  setActiveTab,
  cartCount,
  isFabOpen,
  toggleFabMenu,
}: MobileBottomNavProps) {
  return (
    <>
      {/* BOTTOM NAVIGATION */}
      <nav className="bg-white border-t border-slate-200 pb-safe z-50 fixed bottom-0 left-0 w-full rounded-t-[2rem] shadow-[0_-10px_40px_rgba(0,0,0,0.08)] md:hidden">
        <div className="flex justify-around items-center h-16 px-2 relative">
          <button
            onClick={() => setActiveTab("categories")}
            className={`flex flex-col items-center gap-1 w-16 transition-colors ${activeTab === "categories" ? "text-orange-600" : "text-slate-400 hover:text-slate-600"}`}
          >
            <Grid
              size={22}
              className={
                activeTab === "categories"
                  ? "fill-orange-50 text-orange-600"
                  : ""
              }
            />
            <span
              className={`text-[10px] ${activeTab === "categories" ? "font-bold" : "font-medium"}`}
            >
              Danh mục
            </span>
          </button>

          <Link
            href="/quick-order"
            className={`flex flex-col items-center justify-center gap-1 w-16 transition-colors ${activeTab === "home" ? "text-orange-600" : "text-slate-400 hover:text-slate-600"}`}
          >
            <ShoppingBag
              size={22}
              className={
                activeTab === "home" ? "fill-orange-50 text-orange-600" : ""
              }
            />
            <span
              className={`text-[10px] ${activeTab === "home" ? "font-bold" : "font-medium"}`}
            >
              Đặt hàng
            </span>
          </Link>

          <div className="w-16"></div>

          <Link
            href="/cart"
            className={`flex flex-col items-center justify-center gap-1 w-16 transition-colors relative ${activeTab === "cart" ? "text-orange-600" : "text-slate-400 hover:text-slate-600"}`}
          >
            <ShoppingCart
              size={22}
              className={
                activeTab === "cart" ? "fill-orange-50 text-orange-600" : ""
              }
            />
            {cartCount > 0 && (
              <span className="absolute top-0 right-3 bg-red-500 text-white text-[9px] font-black w-3.5 h-3.5 rounded-full flex items-center justify-center ring-2 ring-white">
                {cartCount}
              </span>
            )}
            <span
              className={`text-[10px] ${activeTab === "cart" ? "font-bold" : "font-medium"}`}
            >
              Giỏ hàng
            </span>
          </Link>

          <button
            onClick={() => setActiveTab("profile")}
            className={`flex flex-col items-center gap-1 w-16 transition-colors relative ${activeTab === "profile" ? "text-orange-600" : "text-slate-400 hover:text-slate-600"}`}
          >
            <User
              size={22}
              className={
                activeTab === "profile" ? "fill-orange-50 text-orange-600" : ""
              }
            />
            <span className="absolute top-0.5 right-4 w-2 h-2 bg-red-500 rounded-full border border-white"></span>
            <span
              className={`text-[10px] ${activeTab === "profile" ? "font-bold" : "font-medium"}`}
            >
              Cá nhân
            </span>
          </button>
        </div>
      </nav>

      {/* THE LOGO - FLOATING CENTER FAB */}
      <div className="fixed bottom-6 left-1/2 -translate-x-1/2 z-[60] md:hidden">
        <button
          onClick={toggleFabMenu}
          className={`w-16 h-16 rounded-full flex items-center justify-center bg-white shadow-[0_8px_30px_rgba(249,115,22,0.3)] border-4 transition-all duration-500 active:scale-90 overflow-hidden ${
            isFabOpen
              ? "rotate-[360deg] scale-110 border-orange-200  shadow-orange-500/50"
              : "border-slate-50 hover:scale-105 hover:-translate-y-1"
          }`}
        >
          <img
            src="/logo.png"
            alt="Nam Viet"
            className="w-full h-full object-cover"
          />
        </button>
      </div>
    </>
  );
}

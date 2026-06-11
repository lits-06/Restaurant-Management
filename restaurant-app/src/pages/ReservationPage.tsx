import React, { useState, useMemo, useRef, useEffect } from 'react';
import { menuApi, tableApi, ordersApi, type MenuItemDto } from '../api/gateway.api';
import { useAuthStore } from '../store/authStore';

interface ReservationPageProps {
  onNeedLogin?: () => void;
}

// Định nghĩa cấu trúc dữ liệu cho món ăn
interface Dish {
  id: string;
  name: string;
  price: number;
  description: string;
  category: 'Appetizers' | 'Main Courses';
  image: string;
}

// Mảng dữ liệu mẫu (Sau này bạn có thể thay thế bằng dữ liệu gọi từ API)
const mockDishes: Dish[] = [
  {
    id: 'app-1',
    name: 'Black Truffle Carpaccio',
    price: 28,
    description: 'Prime Wagyu beef, shaved Perigord truffles.',
    category: 'Appetizers',
    image: 'https://www.gstatic.com/labs-code/stitch/stitch-placeholder-300x300.svg',
  },
  {
    id: 'app-2',
    name: 'Heirloom Tomato & Burrata',
    price: 24,
    description: 'Creamy Puglia burrata, balsamic reduction.',
    category: 'Appetizers',
    image: 'https://www.gstatic.com/labs-code/stitch/stitch-placeholder-300x300.svg',
  },
  {
    id: 'main-1',
    name: 'Herb-Crusted Rack of Lamb',
    price: 54,
    description: 'Provencal herb crust, mint pea purée.',
    category: 'Main Courses',
    image: 'https://www.gstatic.com/labs-code/stitch/stitch-placeholder-300x300.svg',
  },
  {
    id: 'main-2',
    name: 'Miso Glazed Black Cod',
    price: 48,
    description: 'Sustainably sourced cod, ginger-infused dashi.',
    category: 'Main Courses',
    image: 'https://www.gstatic.com/labs-code/stitch/stitch-placeholder-300x300.svg',
  },
];

const getMenuItemId = (item: MenuItemDto) => item.item_id ?? item.itemId ?? '';
const getMenuItemImage = (item: MenuItemDto) => {
  const image = item.image_url ?? item.imageUrl;

  if (!image) {
    return '/images/default-food.jpg';
  }

  // nếu DB trả về assets/images/...
  if (image.startsWith('assets/')) {
    return image.replace('assets/', '/');
  }

  return image.startsWith('/') ? image : `/${image}`;
};

const addHours = (timeStr: string, hours: number): string => {
  const [h, m] = timeStr.split(':').map(Number);
  const totalMin = h * 60 + m + hours * 60;
  const newH = Math.floor(totalMin / 60) % 24;
  const newM = totalMin % 60;
  return `${String(newH).padStart(2, '0')}:${String(newM).padStart(2, '0')}`;
};

const toDish = (item: MenuItemDto): Dish => ({
  id: getMenuItemId(item),
  name: item.name ?? 'Unnamed dish',
  price: item.price ?? 0,
  description: item.description ?? '',
  category: item.category?.toLowerCase().includes('app') || item.category?.toLowerCase().includes('khai')
    ? 'Appetizers'
    : 'Main Courses',
  image: getMenuItemImage(item),
});

const ReservationPage: React.FC<ReservationPageProps> = ({ onNeedLogin }) => {
  // Quản lý các bước (Step 1, 2, 3)
  const [step, setStep] = useState<number>(1);

  // Quản lý số lượng món ăn đã chọn: { [dishId]: quantity }
  // Khởi tạo sẵn giá trị giống bản thiết kế cũ (Tomato x1, Lamb x2)
  const [cart, setCart] = useState<Record<string, number>>({
  });
  const [dishes, setDishes] = useState<Dish[]>(mockDishes);
  const [menuError, setMenuError] = useState<string>('');
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [submitError, setSubmitError] = useState<string>('');
  const [confirmedOrderId, setConfirmedOrderId] = useState<string>('');
  const [confirmedTableNumber, setConfirmedTableNumber] = useState<number | null>(null);

  // Thông tin đặt bàn cơ bản
  const [bookingDate, setBookingDate] = useState<string>(() => new Date().toISOString().slice(0, 10));
  const [guestCount, setGuestCount] = useState<string>('2 Guests');
  const [selectedTime, setSelectedTime] = useState<string>('19:00');
  const [selectedEndTime, setSelectedEndTime] = useState<string>(() => addHours('19:00', 2));

  const { user } = useAuthStore();

  const [customerName, setCustomerName] = useState<string>('');
  const [phoneNumber, setPhoneNumber] = useState<string>('');
  const [specialNotes, setSpecialNotes] = useState<string>('');

  const referenceNumber = useMemo(() => confirmedOrderId || 'Pending confirmation', [confirmedOrderId]);

  // Canvas Reference phục vụ hiệu ứng Confetti
  const canvasRef = useRef<HTMLCanvasElement | null>(null);

  // Xử lý tăng/giảm số lượng món ăn
  const updateQuantity = (id: string, delta: number) => {
    setCart((prev) => {
      const currentQty = prev[id] || 0;
      const newQty = currentQty + delta;
      if (newQty <= 0) {
        const { [id]: _, ...rest } = prev;
        return rest;
      }
      return { ...prev, [id]: newQty };
    });
  };

  // Auth guard + pre-fill
  useEffect(() => {
    if (!useAuthStore.getState().user) {
      if (onNeedLogin) onNeedLogin();
      return;
    }
    const u = useAuthStore.getState().user!;
    if (u.full_name) setCustomerName(u.full_name);
    if (u.phone) setPhoneNumber(u.phone);
  }, []);

  useEffect(() => {
    let isMounted = true;

    menuApi
      .listItems({ page: 1, page_size: 100 })
      .then((response) => {
        if (!isMounted) return;

        const apiDishes = (response.items ?? []).map(toDish).filter((dish) => dish.id);
        if (apiDishes.length > 0) {
          setDishes(apiDishes);
          setCart({});
          setMenuError('');
        }
      })
      .catch((error) => {
        if (isMounted) {
          setMenuError(error instanceof Error ? error.message : 'Cannot load menu from server.');
        }
      });

    return () => {
      isMounted = false;
    };
  }, []);

  // Tính toán hóa đơn dựa vào giỏ hàng thực tế
  const preOrderSubtotal = useMemo(() => {
    return Object.entries(cart).reduce((sum, [id, qty]) => {
      const dish = dishes.find((d) => d.id === id);
      return sum + (dish ? dish.price * qty : 0);
    }, 0);
  }, [cart, dishes]);


  // Chuyển bước đi kèm hiệu ứng cuộn mượt mà trên mobile
  const goToStep = (stepNumber: number) => {
    setStep(stepNumber);
    if (window.innerWidth < 1024) {
      window.scrollTo({ top: 200, behavior: 'smooth' });
    }
  };

  const guestCountNumber = useMemo(() => {
    const parsed = Number.parseInt(guestCount, 10);
    return Number.isFinite(parsed) && parsed > 0 ? parsed : 8;
  }, [guestCount]);

  // Khi user đổi giờ bắt đầu → tự động reset giờ kết thúc về start+2h
  useEffect(() => {
    setSelectedEndTime(addHours(selectedTime, 2));
  }, [selectedTime]);

  const handleProceedToPayment = () => {
    goToStep(2);
  };

  const submitOrder = async () => {
    setSubmitError('');

    if (!customerName.trim() || !phoneNumber.trim()) {
      setSubmitError('Please enter your name and phone number before completing.');
      return;
    }

    if (selectedEndTime <= selectedTime) {
      setSubmitError('End time must be after start time.');
      return;
    }

    setIsSubmitting(true);
    try {
      const items = Object.entries(cart).map(([itemId, quantity]) => ({ item_id: itemId, quantity }));

      const response = await ordersApi.create({
        name: customerName.trim(),
        phone: phoneNumber.trim(),
        date: bookingDate,
        time: selectedTime,
        end_time: selectedEndTime,
        party_size: guestCountNumber,
        notes: specialNotes.trim(),
        items,
      });

      setConfirmedOrderId(response.order?.order_id ?? '');

      const tableId = response.order?.table_id;
      if (tableId) {
        tableApi.getOne(tableId)
          .then((r) => setConfirmedTableNumber(r.table?.table_number ?? null))
          .catch(() => {});
      }

      goToStep(3);
    } catch (error) {
      setSubmitError(error instanceof Error ? error.message : 'Unable to submit reservation to server.');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Hàm helper để lọc món ăn theo danh mục
  const renderDishGroup = (category: 'Appetizers' | 'Main Courses') => {
    return dishes
      .filter((dish) => dish.category === category)
      .map((dish) => {
        const quantity = cart[dish.id] || 0;
        const isSelected = quantity > 0;

        return (
          <div
            key={dish.id}
            className={`group flex border rounded-xl overflow-hidden transition-all ${
              isSelected
                ? 'border-2 border-primary bg-primary/5 shadow-sm'
                : 'border-outline-variant bg-surface hover:shadow-sm'
            }`}
          >
            <div className="w-24 h-24 sm:w-32 sm:h-auto flex-shrink-0">
              <img alt={dish.name} className="w-full h-full object-cover" src={dish.image} />
            </div>
            <div className="p-4 flex-1 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
              <div className="flex-1">
                <div className="flex justify-between items-start">
                  <h4 className="font-bold text-sm">{dish.name}</h4>
                  <span className="text-primary font-bold text-xs ml-2">{dish.price.toLocaleString('vi-VN')} ₫</span>
                </div>
                <p className="text-[10px] text-on-surface-variant leading-tight mt-1 line-clamp-2">
                  {dish.description}
                </p>
                {isSelected && (
                  <span className="inline-flex items-center gap-1 text-[10px] font-bold text-primary mt-1">
                    <span
                      className="material-symbols-outlined text-xs"
                      style={{ fontVariationSettings: "'FILL' 1" }}
                    >
                      check_circle
                    </span>{' '}
                    Selected
                  </span>
                )}
              </div>
              <div
                className={`flex items-center gap-3 rounded-full px-2 py-1 border ${
                  isSelected ? 'bg-white border-primary' : 'bg-surface-container-low border-outline-variant/30'
                }`}
              >
                <button
                  type="button"
                  onClick={() => updateQuantity(dish.id, -1)}
                  className="w-7 h-7 flex items-center justify-center rounded-full hover:bg-surface-container-high text-primary"
                >
                  <span className="material-symbols-outlined text-lg">remove</span>
                </button>
                <span className="font-bold text-sm w-4 text-center">{quantity}</span>
                <button
                  type="button"
                  onClick={() => updateQuantity(dish.id, 1)}
                  className="w-7 h-7 flex items-center justify-center rounded-full bg-primary text-on-primary hover:opacity-90"
                >
                  <span className="material-symbols-outlined text-lg">add</span>
                </button>
              </div>
            </div>
          </div>
        );
      });
  };

  // Hiệu ứng pháo hoa giấy (Confetti) khi vào Bước 3
  useEffect(() => {
    if (step !== 3 || !canvasRef.current) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    let animationFrameId: number;
    
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;

    let particles: any[] = [];
    const colors = ['#735c00', '#d4af37', '#ffe088', '#e9c349'];

    class Particle {
      x: number;
      y: number;
      size: number;
      speedX: number;
      speedY: number;
      color: string;
      rotation: number;
      rotationSpeed: number;

      constructor() {
        this.x = Math.random() * canvas.width;
        this.y = Math.random() * canvas.height - canvas.height;
        this.size = Math.random() * 8 + 4;
        this.speedX = Math.random() * 3 - 1.5;
        this.speedY = Math.random() * 3 + 2;
        this.color = colors[Math.floor(Math.random() * colors.length)];
        this.rotation = Math.random() * 360;
        this.rotationSpeed = Math.random() * 10 - 5;
      }
      update() {
        this.y += this.speedY;
        this.x += this.speedX;
        this.rotation += this.rotationSpeed;
      }
      draw() {
        if (!ctx) return;
        ctx.save();
        ctx.translate(this.x, this.y);
        ctx.rotate((this.rotation * Math.PI) / 180);
        ctx.fillStyle = this.color;
        ctx.fillRect(-this.size / 2, -this.size / 2, this.size, this.size);
        ctx.restore();
      }
    }

    const initConfetti = () => {
      for (let i = 0; i < 100; i++) {
        particles.push(new Particle());
      }
    };

    const animateConfetti = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      particles.forEach((p, index) => {
        p.update();
        p.draw();
        if (p.y > canvas.height) {
          particles.splice(index, 1);
        }
      });
      if (particles.length > 0) {
        animationFrameId = requestAnimationFrame(animateConfetti);
      }
    };

    const handleResize = () => {
      canvas.width = window.innerWidth;
      canvas.height = window.innerHeight;
    };

    initConfetti();
    animateConfetti();

    window.addEventListener('resize', handleResize);

    // Dọn dẹp sự kiện và vòng lặp khi chuyển bước hoặc unmount Component
    return () => {
      cancelAnimationFrame(animationFrameId);
      window.removeEventListener('resize', handleResize);
    };
  }, [step]);

  // Định dạng ngày hiển thị (VD: Friday, Oct 24)
  const formattedDate = useMemo(() => {
    try {
      const options: Intl.DateTimeFormatOptions = { weekday: 'long', month: 'short', day: 'numeric' };
      return new Date(bookingDate).toLocaleDateString('en-US', options);
    } catch {
      return 'Friday, Oct 24';
    }
  }, [bookingDate]);

  return (
    <div className="bg-background text-on-surface font-body-md overflow-x-hidden min-h-screen flex flex-col">
      {/* Canvas phục vụ hiệu ứng Confetti */}
      {step === 3 && (
        <canvas
          ref={canvasRef}
          className="fixed top-0 left-0 w-full h-full pointer-events-none z-50"
        />
      )}

      <main className="max-w-container-max mx-auto px-margin-mobile md:px-margin-desktop py-12 flex-1 w-full">
        {/* GIẢI QUYẾT VẤN ĐỀ 1: Page Header chỉ xuất hiện ở bước 1 */}
        {step === 1 && (
          <div className="mb-12 animate-in fade-in duration-300">
            <h1 className="font-headline-xl text-headline-xl text-on-surface mb-2">Secure Your Table</h1>
            <p className="font-body-lg text-body-lg text-on-surface-variant max-w-2xl">
              Experience culinary artistry at its peak. Our booking process ensures a seamless transition from anticipation to indulgence.
            </p>
          </div>
        )}

        {/* THAY ĐỔI 1: Đưa Progress Stepper lên trên grid để trải rộng, thoáng đãng hơn */}
        <div className="w-full max-w-5xl mx-auto flex items-center justify-between mb-12 px-4">
          {[
            { number: 1, label: 'Details' },
            { number: 2, label: 'Payment' },
            { number: 3, label: 'Confirm' },
          ].map((node, index) => {
            const isActive = step === node.number;
            const isCompleted = step > node.number;

            return (
              <React.Fragment key={node.number}>
                <div className="flex flex-col items-center gap-2">
                  <div
                    className={`w-10 h-10 rounded-full flex items-center justify-center font-bold transition-all ${
                      isActive
                        ? 'bg-primary text-on-primary shadow-md scale-110 border-2 border-primary'
                        : isCompleted
                        ? 'bg-primary text-on-primary shadow-sm'
                        : 'bg-surface-container-high text-on-surface-variant'
                    }`}
                  >
                    {isCompleted ? (
                      <span className="material-symbols-outlined text-xl">check</span>
                    ) : (
                      node.number
                    )}
                  </div>
                  <span
                    className={`text-label-sm font-label-sm ${
                      isActive ? 'text-primary font-bold' : 'text-on-surface-variant'
                    }`}
                  >
                    {node.label}
                  </span>
                </div>
                {index < 2 && (
                  <div 
                    className={`flex-1 h-[2px] mx-6 mb-6 transition-all ${
                      step > node.number ? 'bg-primary' : 'bg-outline-variant'
                    }`} 
                  />
                )}
              </React.Fragment>
            );
          })}
        </div>

        <div className="flex flex-col lg:grid lg:grid-cols-12 gap-gutter">
          {/* Left: Booking Form */}
          <div className={step === 3 ? "lg:col-span-12" : "lg:col-span-8"}>

            {/* Form Card */}
            <div className="bg-surface-container-lowest rounded-xl shadow-[0_4px_12px_rgba(0,0,0,0.04)] p-8 border border-outline-variant/30">
              
              {/* Step 1: Selection */}
              {step === 1 && (
                <div className="transition-all duration-300">
                  <h2 className="font-headline-md text-headline-md mb-8">Reservation Particulars</h2>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-10">
                    {/* Date */}
                    <div className="space-y-2">
                      <label className="text-label-sm font-label-sm text-on-surface-variant block uppercase">Preferred Date</label>
                      <input
                        className="w-full bg-surface border border-outline-variant rounded-lg px-4 py-3 focus:border-primary focus:ring-1 focus:ring-primary/10 transition-all outline-none"
                        type="date"
                        value={bookingDate}
                        onChange={(e) => setBookingDate(e.target.value)}
                      />
                    </div>
                    {/* Guests */}
                    <div className="space-y-2">
                      <label className="text-label-sm font-label-sm text-on-surface-variant block uppercase">Guest Count</label>
                      <select
                        className="w-full bg-surface border border-outline-variant rounded-lg px-4 py-3 focus:border-primary focus:ring-1 focus:ring-primary/10 transition-all outline-none"
                        value={guestCount}
                        onChange={(e) => setGuestCount(e.target.value)}
                      >
                        <option>2 Guests</option>
                        <option>4 Guests</option>
                        <option>6 Guests</option>
                        <option>Private Party (8+)</option>
                      </select>
                    </div>
                    {/* Start time */}
                    <div className="space-y-2">
                      <label className="text-label-sm font-label-sm text-on-surface-variant block uppercase">
                        Start Time
                      </label>
                      <input
                        className="w-full bg-surface border border-outline-variant rounded-lg px-4 py-3 focus:border-primary focus:ring-1 focus:ring-primary/10 transition-all outline-none"
                        type="time"
                        value={selectedTime}
                        min="10:00"
                        max="22:00"
                        onChange={(e) => setSelectedTime(e.target.value)}
                      />
                      <p className="text-xs text-on-surface-variant">Restaurant open 10:00 – 22:00</p>
                    </div>
                    {/* End time */}
                    <div className="space-y-2">
                      <label className="text-label-sm font-label-sm text-on-surface-variant block uppercase">
                        End Time
                      </label>
                      <input
                        className={`w-full bg-surface border rounded-lg px-4 py-3 focus:ring-1 transition-all outline-none ${
                          selectedEndTime <= selectedTime
                            ? 'border-error focus:border-error focus:ring-error/10'
                            : 'border-outline-variant focus:border-primary focus:ring-primary/10'
                        }`}
                        type="time"
                        value={selectedEndTime}
                        min={selectedTime}
                        max="22:00"
                        onChange={(e) => setSelectedEndTime(e.target.value)}
                      />
                      {selectedEndTime <= selectedTime
                        ? <p className="text-xs text-error">Must be after start time</p>
                        : <p className="text-xs text-on-surface-variant">Default +2h, adjust if needed</p>
                      }
                    </div>
                  </div>

                  {/* Pre-order Selection Feature */}
                  <div className="border-t border-outline-variant/30 pt-10 space-y-8">
                    <div className="flex flex-col gap-2">
                      <h2 className="font-headline-md text-headline-md">Pre-order Selection</h2>
                      <p className="text-body-md text-on-surface-variant">
                        Enhance your evening by pre-selecting our signature dishes. This ensures priority preparation for your arrival.
                      </p>
                      {menuError && (
                        <p className="text-xs text-error">
                          Menu API unavailable, showing sample data: {menuError}
                        </p>
                      )}
                    </div>
                    <div className="max-h-[600px] overflow-y-auto pr-2 space-y-10 custom-scrollbar">
                      {/* Appetizers */}
                      <div className="space-y-4">
                        <h3 className="font-label-sm text-primary uppercase tracking-widest border-b border-outline-variant pb-2 sticky top-0 bg-surface-container-lowest z-10">
                          Appetizers
                        </h3>
                        <div className="grid grid-cols-1 gap-4">{renderDishGroup('Appetizers')}</div>
                      </div>
                      {/* Main Courses */}
                      <div className="space-y-4">
                        <h3 className="font-label-sm text-primary uppercase tracking-widest border-b border-outline-variant pb-2 sticky top-0 bg-surface-container-lowest z-10">
                          Main Courses
                        </h3>
                        <div className="grid grid-cols-1 gap-4">{renderDishGroup('Main Courses')}</div>
                      </div>
                    </div>
                  </div>

                  <div className="mt-10 flex justify-end">
                    <button
                      type="button"
                      onClick={handleProceedToPayment}
                      className="bg-primary text-on-primary px-10 py-4 rounded-lg font-bold hover:shadow-lg transition-all active:scale-95 flex items-center gap-2"
                    >
                      Proceed to Payment
                      <span className="material-symbols-outlined">arrow_forward</span>
                    </button>
                  </div>
                </div>
              )}

              {/* Step 2: Payment */}
              {step === 2 && (
                <div className="transition-all duration-300">
                  {/* Thông tin khách hàng */}
                  <section className="mb-12 border-b border-outline-variant/30 pb-10">
                    <h2 className="font-headline-lg text-headline-lg mb-8 text-xl font-bold">Guest Information</h2>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div className="space-y-2">
                        <label className="block text-sm font-medium text-on-surface" htmlFor="full-name">Full Name</label>
                        <input
                          className="w-full px-4 py-3 bg-surface border border-outline-variant rounded-lg focus:ring-1 focus:ring-primary focus:border-primary font-body-md outline-none transition-all"
                          id="full-name"
                          placeholder="John Doe"
                          type="text"
                          value={customerName}
                          onChange={(e) => setCustomerName(e.target.value)}
                        />
                      </div>
                      <div className="space-y-2">
                        <label className="block text-sm font-medium text-on-surface" htmlFor="phone-number">Phone Number</label>
                        <input
                          className="w-full px-4 py-3 bg-surface border border-outline-variant rounded-lg focus:ring-1 focus:ring-primary focus:border-primary font-body-md outline-none transition-all"
                          id="phone-number"
                          placeholder="+1 555 123 4567"
                          type="tel"
                          value={phoneNumber}
                          onChange={(e) => setPhoneNumber(e.target.value)}
                        />
                      </div>
                      <div className="space-y-2 md:col-span-2">
                        <label className="block text-sm font-medium text-on-surface" htmlFor="special-notes">Special Notes</label>
                        <textarea
                          className="w-full px-4 py-3 bg-surface border border-outline-variant rounded-lg focus:ring-1 focus:ring-primary focus:border-primary font-body-md min-h-[120px] outline-none transition-all"
                          id="special-notes"
                          placeholder="e.g. Seafood allergy, window seat preferred..."
                          value={specialNotes}
                          onChange={(e) => setSpecialNotes(e.target.value)}
                        ></textarea>
                      </div>
                    </div>
                  </section>

                  <div className="py-8 border-t border-outline-variant/30">
                    <div className="flex flex-col items-center text-center space-y-4">
                      <div className="w-16 h-16 rounded-full bg-primary-container/20 flex items-center justify-center text-primary">
                        <span className="material-symbols-outlined text-3xl">restaurant</span>
                      </div>
                      <h2 className="font-headline-md font-bold">Pay at Restaurant</h2>
                      <p className="font-body-lg text-on-surface-variant max-w-md">
                        Payment will be settled directly at the restaurant upon your arrival.
                      </p>
                    </div>
                  </div>

                  {/* Nút điều hướng biểu mẫu */}
                  <div className="flex items-center justify-between pt-8 border-t border-outline-variant/30 mt-8">
                    <button 
                      type="button"
                      onClick={() => goToStep(1)}
                      className="flex items-center gap-2 text-on-surface-variant hover:text-primary transition-all font-label-sm font-medium"
                    >
                      <span className="material-symbols-outlined">arrow_back</span>
                      Back to Reservation Details
                    </button>
                    <button 
                      type="button"
                      onClick={submitOrder}
                      disabled={isSubmitting}
                      className="bg-primary hover:bg-primary/90 text-on-primary px-10 py-4 rounded-lg font-bold shadow-lg transform active:scale-95 transition-all disabled:opacity-60"
                    >
                      {isSubmitting ? 'Sending...' : 'Complete Payment'}
                    </button>
                  </div>
                  {submitError && <p className="mt-4 text-right text-sm text-error">{submitError}</p>}
                </div>
              )}

              {step === 3 && (
                <div className="animate-in fade-in zoom-in duration-500">
                  {/* Success Content Card */}
                  <div className="bg-surface-container-lowest rounded-xl shadow-lg border border-outline-variant overflow-hidden">
                    <div className="p-8 md:p-12 text-center border-b border-outline-variant bg-gradient-to-b from-primary-container/10 to-transparent">
                      <div className="w-20 h-20 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-6">
                        <span className="material-symbols-outlined text-primary text-5xl" style={{ fontVariationSettings: "'wght' 600" }}>
                          check_circle
                        </span>
                      </div>
                      <h1 className="font-headline-lg text-headline-lg text-on-surface mb-2 text-2xl font-bold">Reservation Confirmed</h1>
                      <p className="font-body-md text-on-surface-variant max-w-md mx-auto">
                        We're delighted to host you. A confirmation email has been sent to your inbox.
                      </p>
                    </div>

                    <div className="p-8 md:p-12">
                      <div className="grid md:grid-cols-2 gap-8">
                        {/* Left Column: Reservation Summary */}
                        <div className="space-y-6">
                          <h2 className="font-label-sm text-label-sm uppercase tracking-widest text-primary border-b border-outline-variant pb-2 font-bold">
                            Booking Summary
                          </h2>
                          <div className="grid grid-cols-2 gap-y-4 text-sm">
                            <div>
                              <p className="text-xs text-on-surface-variant">Reference Number</p>
                              <p className="font-bold text-on-surface">{referenceNumber}</p>
                            </div>
                            <div>
                              <p className="text-xs text-on-surface-variant">Guests</p>
                              <p className="font-bold text-on-surface">{guestCount}</p>
                            </div>
                            <div className="col-span-2">
                              <p className="text-xs text-on-surface-variant">Date & Time</p>
                              <p className="font-bold text-on-surface">{formattedDate} • {selectedTime} – {selectedEndTime}</p>
                            </div>
                            <div className="col-span-2">
                              <p className="text-xs text-on-surface-variant">Assigned Table</p>
                              {confirmedTableNumber !== null
                                ? <p className="font-bold text-on-surface">Table {confirmedTableNumber}</p>
                                : <p className="text-on-surface-variant italic text-xs">Staff will assist you upon arrival</p>
                              }
                            </div>
                            {customerName && (
                              <div className="col-span-2">
                                <p className="text-xs text-on-surface-variant">Customer Details</p>
                                <p className="font-bold text-on-surface">{customerName} • {phoneNumber}</p>
                              </div>
                            )}
                          </div>
                        </div>

                        {/* Right Column: Pre-ordered Items */}
                        <div className="space-y-6">
                          <h2 className="font-label-sm text-label-sm uppercase tracking-widest text-primary border-b border-outline-variant pb-2 font-bold">
                            Pre-ordered Selection
                          </h2>
                          <ul className="space-y-3 text-sm">
                            {Object.keys(cart).length > 0 ? (
                              Object.entries(cart).map(([id, qty]) => {
                                const dish = dishes.find((d) => d.id === id);
                                if (!dish) return null;
                                return (
                                  <li key={id} className="flex justify-between items-center">
                                    <span className="text-on-surface">{dish.name}</span>
                                    <span className="text-on-surface-variant font-semibold">x{qty}</span>
                                  </li>
                                );
                              })
                            ) : (
                              <li className="text-on-surface-variant italic">No items pre-ordered</li>
                            )}
                            
                            <li className="flex justify-between items-center pt-4 border-t border-dashed border-outline-variant">
                              <span className="font-bold text-primary">Pre-order Total</span>
                              <span className="font-bold text-primary text-base">{preOrderSubtotal.toLocaleString('vi-VN')} ₫</span>
                            </li>
                          </ul>
                        </div>
                      </div>

                      {/* Action Buttons */}
                      <div className="mt-12 flex flex-col md:flex-row gap-4 justify-center">
                        <button 
                          onClick={() => goToStep(1)}
                          className="flex items-center justify-center gap-2 px-8 py-3 bg-transparent text-primary rounded-lg font-label-sm text-label-sm hover:bg-primary/5 transition-all active:scale-95 font-semibold"
                        >
                          <span className="material-symbols-outlined text-lg">home</span>
                          Return to Home
                        </button>
                      </div>
                    </div>
                  </div>

                  {/* Venue Preview Banner */}
                  <div className="mt-12 rounded-xl overflow-hidden relative h-48 group shadow-lg cursor-pointer">
                    <img 
                      alt="LuxeBistro Dining Room" 
                      className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105" 
                      src="https://lh3.googleusercontent.com/aida-public/AB6AXuAbZta_8YeE0OhA_1I4PhEJ7CTYylmxOVn9oGpu15ZrV3itv7FT6kkqXMJO1nA_oHSdEifWDpI95uRskTEtBtn1brYolsko_Jy2f1JZVbcPVVXmaVQQyiG8DsULh-5dXeKqWSyY3qi9Be1dYDsRowdKtVXaleT7ZsSFiLC3sCoqGqW3Fv0o3_ULVVJfL1CHQ_tpNaxl9w_tlg8wJVK9RfoFKOHbmYuiwA0Zp7RvOFfrjwXZVhi5yoYbafK_TI-4SHE_VWVS6BrcX-M"
                    />
                    <div className="absolute inset-0 bg-gradient-to-t from-black/80 to-transparent flex flex-col justify-end p-6">
                      <p className="text-white font-bold text-xl">See you on Friday</p>
                      <p className="text-white/80 text-sm">123 Culinary Way, Gastronomy District</p>
                    </div>
                  </div>
                </div>
              )}

            </div>
          </div>

          {/* THAY ĐỔI 2: Ẩn hoàn toàn Sidebar khi ở Step 3 */}
          {step !== 3 && (
            <aside className="lg:col-span-4 space-y-gutter">
              <div className="bg-surface-container-lowest rounded-xl shadow-md border border-outline-variant/30 overflow-hidden sticky top-24">
                <div className="h-32 bg-primary-container/20 relative">
                  <img
                    alt="Restaurant Interior"
                    className="w-full h-full object-cover opacity-60"
                    src="https://lh3.googleusercontent.com/aida-public/AB6AXuAwxEpxbCg8RCoiEAJeC6QpNwyrLqG5W4IqVJRHjAgDjHPxziJ-dw3zZOVWzj5vieZD0LH5Q6hnAOy0lmkuHeharUOZIBMA_4w6kxuwyT85D6ARAl7bg2F9ZJT13OsxoezbD_OI_o9swMO3YNl8CAihKJc-Nm4OMEZgi6h-3CWVS25eP5Eyaw72hdnkjpkPLfjA-dc-r0MnUrkTkwB6qUxT7HgZi2DNdXdQoKjuw5Gc2v0ddLiPX9qh7NWkppTihp-znLiazavMneQ"
                  />
                  <div className="absolute bottom-4 left-6">
                    <span className="bg-primary text-on-primary text-[10px] font-bold px-2 py-1 rounded uppercase tracking-widest">
                      Selected Experience
                    </span>
                  </div>
                </div>
                <div className="p-6 space-y-6">
                  <h3 className="font-headline-md text-headline-md border-b border-outline-variant pb-4">
                    Reservation Summary
                  </h3>
                  <div className="space-y-4">
                    <div className="flex justify-between items-center">
                      <span className="text-on-surface-variant font-medium">Restaurant</span>
                      <span className="font-bold text-on-surface">LuxeBistro</span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-on-surface-variant font-medium">Date &amp; Time</span>
                      <span className="font-bold text-on-surface">
                        {bookingDate}, {selectedTime} – {selectedEndTime}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-on-surface-variant font-medium">Guests</span>
                      <span className="font-bold text-on-surface">{guestCount}</span>
                    </div>

                    {/* Pre-ordered Items Section trong Sidebar */}
                    {Object.keys(cart).length > 0 && (
                      <div className="pt-2">
                        <span className="text-on-surface-variant font-medium block mb-2">Pre-ordered Items</span>
                        <div className="space-y-2 bg-surface-container-low p-3 rounded-lg border border-outline-variant/30">
                          {Object.entries(cart).map(([id, qty]) => {
                            const dish = dishes.find((d) => d.id === id);
                            if (!dish) return null;
                            return (
                              <div key={id} className="flex justify-between text-sm">
                                <span className="text-on-surface">
                                  {dish.name} <span className="text-xs font-bold text-primary ml-1">x{qty}</span>
                                </span>
                                <span className="font-bold">{(dish.price * qty).toLocaleString('vi-VN')} ₫</span>
                              </div>
                            );
                          })}
                        </div>
                      </div>
                    )}
                  </div>

                  <div className="pt-6 border-t border-outline-variant">
                    <div className="flex justify-between items-center">
                      <span className="font-headline-md text-headline-md text-primary">Pre-order Total</span>
                      <span className="font-headline-md text-headline-md text-primary">{preOrderSubtotal.toLocaleString('vi-VN')} ₫</span>
                    </div>
                  </div>

                  <div className="bg-primary-container/10 p-4 rounded-lg flex items-start gap-3">
                    <span className="material-symbols-outlined text-primary text-xl">info</span>
                    <p className="text-label-sm text-on-primary-container">
                      Payment settled at the restaurant. Total shown is pre-order estimate only.
                    </p>
                  </div>
                </div>
              </div>

              {/* Assistance Card */}
              <div className="bg-inverse-surface text-inverse-on-surface p-6 rounded-xl">
                <h4 className="font-bold mb-2">Need Assistance?</h4>
                <p className="text-label-sm opacity-80 mb-4">
                  Our concierge is available 24/7 to help with special requests or large group arrangements.
                </p>
                <button type="button" className="w-full py-2 rounded border border-inverse-on-surface/30 hover:bg-inverse-on-surface/10 transition-colors font-bold text-sm">
                  Call Concierge
                </button>
              </div>
            </aside>
          )}


        </div>
      </main>
    </div>
  );
};

export default ReservationPage;

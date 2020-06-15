#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/init.h>
#include <linux/sched/signal.h>
#include <linux/sched.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h> 
#include <asm/uaccess.h> 
#include <linux/hugetlb.h>
#include <linux/init.h>
#include <linux/fs.h>
#include <linux/string.h> 


MODULE_LICENSE("GPL");
MODULE_AUTHOR("201503935 - Douglas Daniel Aguilar Cuque");
MODULE_DESCRIPTION("Obtener los procesos y sus hijos y escribir un archivo");

#define ENTRY_NAME "cpu_201503935"
#define PERMS 0644 
#define PARENT NULL 
#define BUFSIZE  	150

struct task_struct *task;        
struct task_struct *task_child;        
struct list_head *list;            
    

static int escribirArchivo(struct seq_file * archivo, void *v) {	
    char estado[20];
    char mensaje[150];
    seq_printf(archivo, "Carnet: 201503935\n");
    seq_printf(archivo, "Nombre: Douglas Daniel Aguilar Cuque\n");
    for_each_process(task){           
        if(task->state == 0){
            strcpy(estado, "Running");
        }else if(task->state == 1){
            strcpy(estado, "Interruptible");
        }else if(task->state == 2){
            strcpy(estado, "Uninterruptible");
        }else if(task->state == 3){
            strcpy(estado, "Zombie");
        } else if(task->state == 4){
            strcpy(estado, "Stopped");
        }else if(task->state == 5){
            strcpy(estado, "Swapping");
        }else{
            strcpy(estado, "Stopped");
        }
        strcpy(mensaje,"\nPID: %d, Nombre: %s, User: %d, Estado: ");
        strcat(mensaje,estado);
        seq_printf(archivo, mensaje, task->pid, task->comm, task->cred->uid);
        list_for_each(list, &task->children){              
 
            task_child = list_entry( list, struct task_struct, sibling );   
            if(task_child->state == 0){
                strcpy(estado, "Running");
            }else if(task_child->state == 1){
                strcpy(estado, "Interruptible");
            }else if(task_child->state == 2){
                strcpy(estado, "Uninterruptible");
            }else if(task_child->state == 3){
                strcpy(estado, "Zombie");
            } else if(task_child->state == 4){
                strcpy(estado, "Stopped");
            }else if(task_child->state == 5){
                strcpy(estado, "Swapping");
            }else{
                strcpy(estado, "Stopped");
            }
            strcpy(mensaje,"\nHijo de %s,%d, PID: %d, Nombre: %s, User: %d, Estado: ");
            strcat(mensaje, estado);
            seq_printf(archivo, mensaje,task->comm, task->pid, task_child->pid, task_child->comm, task_child->cred->uid);
        }
    }    
    return 0;
}

static int abrir(struct inode *inode, struct  file *file) {
  return single_open(file, escribirArchivo, NULL);
}

static struct file_operations ope =
{    
    .open = abrir,
    .read = seq_read
};

static int inicio(void)
{
    printk(KERN_INFO "Carnet: 201503935\n");
    if(!proc_create(ENTRY_NAME, PERMS, NULL, &ope))
    {
        printk("Error al crear ");
        remove_proc_entry(ENTRY_NAME,NULL);
        return -ENOMEM;
    }
    return 0;
}
 
static void fin(void)
{
    remove_proc_entry(ENTRY_NAME, NULL);
    printk(KERN_INFO "Curso: Sistemas Operativos 1\n");
}
 
module_init(inicio);
module_exit(fin);
 
